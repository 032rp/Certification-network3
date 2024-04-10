package main

import (
	"cert/contracts"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	certContract := new(contracts.CertContract)
	studentContarct := new(contracts.StudentContract)

	chaincode, err := contractapi.NewChaincode(certContract, studentContarct)

	if err != nil {
		log.Panicf("Could not create chaincode : %v", err)
	}

	err = chaincode.Start()

	if err != nil {
		log.Panicf("Failed to start chaincode : %v", err)
	}
}
