package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type Student struct {
	StudentID   string `json:"studentID"`
	Name        string `json:"name"`
	Gender      string `json:"gender"`
	Email       string `json:"email"`
	DateofBirth string `json:"dateofbirth"`
	College     string `json:"college"`
}

type Certificate struct {
	StudentID       string `json:"studentID"`
	CourseID        string `json:"courseID"`
	CertificateHash string `json:"originalcertificatehash"`
	Grade           string `json:"gradereceived"`
}
type Certverify struct {
	CertificateHash string `json:"certificateHash"`
	StudentID       string `json:"studentID"`
	CourseID        string `json:"courseID"`
}

func main() {
	router := gin.Default()

	var wg sync.WaitGroup
	wg.Add(3)
	go ChaincodeEventListener("iit", "mychannel", "CertNetwork", &wg)
	go blockEventListener("iit", "mychannel")
	go pvtBlockEventListener("iit", "mychannel")

	router.GET("/api/event", func(ctx *gin.Context) {
		result := getEvents()
		fmt.Println("result:", result)

		ctx.JSON(http.StatusOK, gin.H{"certificateEvent": result})

	})

	router.POST("/api/registerstudent", func(ctx *gin.Context) {
		var req Student
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
			return
		}

		fmt.Printf("register  %s", req)

		privateData := map[string][]byte{
			"name":        []byte(req.Name),
			"gender":      []byte(req.Gender),
			"email":       []byte(req.Email),
			"dateofbirth": []byte(req.DateofBirth),
			"college":     []byte(req.College),
		}

		submitTxnFn("iit", "mychannel", "CertNetwork", "StudentContract", "private", privateData, "RegisterStudent", req.StudentID)

		ctx.JSON(http.StatusOK, req)
	})

	router.GET("/api/getstudent/:id", func(ctx *gin.Context) {
		studentID := ctx.Param("id")

		result := submitTxnFn("iit", "mychannel", "CertNetwork", "StudentContract", "query", make(map[string][]byte), "GetStudent", studentID)

		ctx.JSON(http.StatusOK, gin.H{"data": result})
	})

	router.POST("/api/issuecertificate", func(ctx *gin.Context) {
		var req Certificate
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
			return
		}

		fmt.Printf("order  %s", req)

		privateData := map[string][]byte{
			// "studentID":               []byte(req.StudentID),
			"originalcertificatehash": []byte(req.CertificateHash),
			"gradereceived":           []byte(req.Grade),
		}

		submitTxnFn("iit", "mychannel", "CertNetwork", "StudentContract", "private", privateData, "IssueCertificate", req.StudentID, req.CourseID)

		ctx.JSON(http.StatusOK, req)
	})

	router.POST("/api/verifycertificate", func(ctx *gin.Context) {
		var req Certverify
		// var certificateHash string
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
			return
		}

		fmt.Printf("certificate response %s", req)
		submitTxnFn("iit", "mychannel", "CertNetwork", "CertContract", "invoke", make(map[string][]byte), "VerifyCertificate", req.CertificateHash, req.StudentID, req.CourseID)

		ctx.JSON(http.StatusOK, req)
	})

	router.GET("/api/getcertificate/:id", func(ctx *gin.Context) {
		courseID := ctx.Param("id")

		result := submitTxnFn("iit", "mychannel", "CertNetwork", "StudentContract", "query", make(map[string][]byte), "GetCertificate", courseID)

		ctx.JSON(http.StatusOK, gin.H{"data": result})
	})

	router.GET("/api/getallstudents", func(ctx *gin.Context) {
		result := submitTxnFn("iit", "mychannel", "CertNetwork", "StudentContract", "query", make(map[string][]byte), "GetAllStudents")
		var students []Student
		if len(result) > 0 {
			// Unmarshal the JSON array string into the student slice
			if err := json.Unmarshal([]byte(result), &students); err != nil {
				fmt.Println("Error:", err)
				return
			}
		}

		ctx.JSON(http.StatusOK, gin.H{
			"title": "All Students ", "studentList": students})
	})

	router.Run("localhost:8080")
}
