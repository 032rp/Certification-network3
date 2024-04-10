package contracts

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type CertContract struct {
	contractapi.Contract
}

type EventData struct {
	Type               string
	StudentID          string
	CourseID           string
	Verifier           string
	VerificationResult string
}

type HistoryQueryResult struct {
	Record    *Certificate `json:"record"`
	TxId      string       `json:"txId"`
	Timestamp string       `json:"timestamp"`
	IsDelete  bool         `json:"isDelete"`
}

func (c *CertContract) VerifyCertificate(ctx contractapi.TransactionContextInterface, certificateHash string, studentID string, courseID string) (*EventData, error) {
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	bytes, err := ctx.GetStub().GetPrivateData(collectionName, courseID)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	if bytes == nil {
		return nil, fmt.Errorf("certificate does not exist with this id %s", courseID)
	}
	var certData Certificate
	err = json.Unmarshal(bytes, &certData)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	if certData.CertificateHash != certificateHash {

		verificationData := EventData{
			Type:               "Certificate Verification",
			StudentID:          studentID,
			CourseID:           courseID,
			Verifier:           clientID,
			VerificationResult: "INVALID CERTIFICATE",
		}

		eventDataByte, _ := json.Marshal(verificationData)
		ctx.GetStub().SetEvent("VerifyCertificate", eventDataByte)

		// return &verificationData, fmt.Errorf("this is not the valid certificate")
		return &verificationData, nil
	}

	verificationData := EventData{
		Type:               "Certificate Verification",
		StudentID:          studentID,
		CourseID:           courseID,
		Verifier:           clientID,
		VerificationResult: "VALID CERTIFICATE",
	}

	eventDataByte, _ := json.Marshal(verificationData)
	ctx.GetStub().SetEvent("VerifyCertificate", eventDataByte)

	return &verificationData, nil
}

func (c *CertContract) GetCertificateHistory(ctx contractapi.TransactionContextInterface, courseID string) ([]*HistoryQueryResult, error) {

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(courseID)
	if err != nil {
		return nil, fmt.Errorf("could not get the data. %s", err)
	}
	defer resultsIterator.Close()

	var records []*HistoryQueryResult
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("could not get the value of resultsIterator. %s", err)
		}

		var certificate Certificate
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &certificate)
			if err != nil {
				return nil, err
			}
		} else {
			certificate = Certificate{
				CourseID: courseID,
			}
		}

		timestamp := response.Timestamp.AsTime()

		formattedTime := timestamp.Format(time.RFC1123)

		record := HistoryQueryResult{
			TxId:      response.TxId,
			Timestamp: formattedTime,
			Record:    &certificate,
			IsDelete:  response.IsDelete,
		}
		records = append(records, &record)
	}

	return records, nil
}
