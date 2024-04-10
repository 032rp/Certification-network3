package contracts

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type StudentContract struct {
	contractapi.Contract
}

type Student struct {
	DocType     string `json:"doctype"`
	Name        string `json:"name"`
	Gender      string `json:"gender"`
	StudentID   string `json:"studentID"`
	Email       string `json:"email"`
	DateofBirth string `json:"dateofbirth"`
	College     string `json:"college"`
	// Certificate string `json:"certificate"`
}

type PaginatedQueryResult struct {
	Records             []*Student `json:"records"`
	FetchedRecordsCount int32      `json:"fetchedRecordsCount"`
	Bookmark            string     `json:"bookmark"`
}

type Certificate struct {
	Doctype         string `json:"doctype"`
	StudentID       string `json:"studentID"`
	CourseID        string `json:"courseID"`
	CertificateHash string `json:"originalcertificatehash"`
	Grade           string `json:"gradereceived"`
}

const collectionName string = "StudentCollection"

func (c *StudentContract) StudentExists(ctx contractapi.TransactionContextInterface, studentID string) (bool, error) {
	data, err := ctx.GetStub().GetPrivateDataHash(collectionName, studentID)
	if err != nil {
		return false, fmt.Errorf("failed to read data from world state: %v", err)
	} else {
		return data != nil, nil
	}
}

func (c *StudentContract) RegisterStudent(ctx contractapi.TransactionContextInterface, studentID string) (string, error) {
	clientID, _ := ctx.GetClientIdentity().GetMSPID()
	if clientID == "iitMSP" {
		exists, err := c.StudentExists(ctx, studentID)
		if err != nil {
			return "", fmt.Errorf("%s", err)
		} else if exists {
			return "", fmt.Errorf("this student is already registered")
		} else {
			var student Student
			transientData, err := ctx.GetStub().GetTransient()
			if err != nil {
				return "", fmt.Errorf("%s", err)
			}
			if len(transientData) == 0 {
				return "", fmt.Errorf("please provide the private data of name, studentID, email, schoolname")
			}

			name, exists := transientData["name"]
			if !exists {
				return "", fmt.Errorf("the name was not specified in transient data. Please try again")
			}
			student.Name = string(name)

			gender, exists := transientData["gender"]
			if !exists {
				return "", fmt.Errorf("gender was not provided in transient data. Please try again")
			}
			student.Gender = string(gender)

			email, exists := transientData["email"]
			if !exists {
				return "", fmt.Errorf("the email was not specified in transient data. Please try again")
			}
			student.Email = string(email)

			dateofbirth, exists := transientData["dateofbirth"]
			if !exists {
				return "", fmt.Errorf("the dateofbirth was not specified in transient data. Please try again")
			}
			student.DateofBirth = string(dateofbirth)

			college, exists := transientData["college"]
			if !exists {
				return "", fmt.Errorf("the college was not specified in transient data. Please try again")
			}
			student.College = string(college)

			student.StudentID = studentID
			student.DocType = "Student"

			bytes, _ := json.Marshal(student)
			err = ctx.GetStub().PutPrivateData(collectionName, studentID, bytes)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("student with studentID %v registered successfully", studentID), nil
		}
	} else {
		return "", fmt.Errorf("you are not authorised to register the student")
	}
}

func (c *StudentContract) GetStudent(ctx contractapi.TransactionContextInterface, studentID string) (*Student, error) {

	exists, err := c.StudentExists(ctx, studentID)
	if err != nil {
		return nil, fmt.Errorf("%s", err)

	} else if !exists {
		return nil, fmt.Errorf("student does not exists with studentID: %v", studentID)
	}
	bytes, err := ctx.GetStub().GetPrivateData(collectionName, studentID)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	var student Student
	err = json.Unmarshal(bytes, &student)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	return &student, nil

}

func (c *StudentContract) IssueCertificate(ctx contractapi.TransactionContextInterface, studentID string, courseID string) (string, error) {
	clientID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("%s", err)
	}
	if clientID == "iitMSP" {
		data, err := ctx.GetStub().GetPrivateData(collectionName, studentID)
		if err != nil {
			return "", fmt.Errorf("%s", err)
		}
		if data == nil {
			return "", fmt.Errorf("student does not exists")
		}

		certID, err := ctx.GetStub().GetPrivateData(collectionName, courseID)
		if err != nil {
			return "", fmt.Errorf("%s", err)
		}
		if certID != nil {
			return "", fmt.Errorf("certificate already exist")
		}
		var certificate Certificate
		transientData, err := ctx.GetStub().GetTransient()
		if err != nil {
			return "", fmt.Errorf("%s", err)
		}
		if len(transientData) == 0 {
			return "", fmt.Errorf("please provide the private data of name, studentID, email, schoolname")
		}

		originalcertificatehash, exists := transientData["originalcertificatehash"]
		if !exists {
			return "", fmt.Errorf("the originalcertificatehash was not specified in transient data. Please try again")
		}
		certificate.CertificateHash = string(originalcertificatehash)

		gradereceived, exists := transientData["gradereceived"]
		if !exists {
			return "", fmt.Errorf("the gradereceived was not specified in transient data. Please try again")
		}
		certificate.Grade = string(gradereceived)

		certificate.CourseID = courseID
		certificate.Doctype = "Certificate"

		bytes, _ := json.Marshal(certificate)
		err = ctx.GetStub().PutPrivateData(collectionName, courseID, bytes)
		if err != nil {
			return "", fmt.Errorf("%s", err)
		}
		return fmt.Sprintf("certficate has been issued for the course %v", courseID), nil

	}
	return "", fmt.Errorf("not authorised to issue certificate")
}

func (c *StudentContract) GetCertificate(ctx contractapi.TransactionContextInterface, courseID string) (*Certificate, error) {

	bytes, err := ctx.GetStub().GetPrivateData(collectionName, courseID)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	if bytes == nil {
		return nil, fmt.Errorf("%s", err)
	}

	var certificate Certificate
	err = json.Unmarshal(bytes, &certificate)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	return &certificate, nil

}
func (c *StudentContract) GetStudentByRange(ctx contractapi.TransactionContextInterface, studentIDstart string, studentIDend string) ([]*Student, error) {
	resultsIterator, err := ctx.GetStub().GetPrivateDataByRange(collectionName, studentIDstart, studentIDend)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	defer resultsIterator.Close()
	return iteratorFunc(resultsIterator)
}

func iteratorFunc(resultsIterator shim.StateQueryIteratorInterface) ([]*Student, error) {
	var students []*Student
	for resultsIterator.HasNext() {
		result, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("unable to iterate")
		}
		var student Student
		err = json.Unmarshal(result.Value, &student)
		if err != nil {
			return nil, fmt.Errorf("%s", err)
		}
		students = append(students, &student)
	}
	return students, nil
}

func (c *StudentContract) GetStudentsWithPagination(ctx contractapi.TransactionContextInterface, pageSize int32, bookmark string) (*PaginatedQueryResult, error) {

	queryString := `{"selector":{"doctype":"Student"}}`

	resultsIterator, responseMetadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, pageSize, bookmark)
	if err != nil {
		return nil, fmt.Errorf("could not get the students records. %s", err)
	}
	defer resultsIterator.Close()

	students, err := iteratorFunc(resultsIterator)
	if err != nil {
		return nil, fmt.Errorf("could not return the car records %s", err)
	}

	return &PaginatedQueryResult{
		Records:             students,
		FetchedRecordsCount: responseMetadata.FetchedRecordsCount,
		Bookmark:            responseMetadata.Bookmark,
	}, nil
}

func (c *StudentContract) GetAllStudents(ctx contractapi.TransactionContextInterface) ([]*Student, error) {
	querystring := `{"selector": {"doctype":"Student"}, "sort":[{"gender": "Male"}]}`

	resultsIterator, err := ctx.GetStub().GetPrivateDataQueryResult(collectionName, querystring)
	if err != nil {
		return nil, fmt.Errorf("could not fetch the result %s", err)
	}
	defer resultsIterator.Close()
	return iteratorFunc(resultsIterator)
}
