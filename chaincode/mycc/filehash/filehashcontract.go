package main

import (
    "fmt"
    "log"

    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// FileHashContract defines the smart contract
type FileHashContract struct {
	contractapi.Contract
}

// StoreFileHash stores the hash of a file in the private data collection
func (s *FileHashContract) StoreFileHash(ctx contractapi.TransactionContextInterface, fileID string, fileHash string) error {
	if fileID == "" || fileHash == "" {
		return fmt.Errorf("fileID and fileHash must not be empty")
	}

	err := ctx.GetStub().PutPrivateData("filehashCollection", fileID, []byte(fileHash))
	if err != nil {
		return fmt.Errorf("failed to store file hash in private data collection: %w", err)
	}
	return nil
}

// QueryFileHash retrieves a file hash from the private data collection
func (s *FileHashContract) QueryFileHash(ctx contractapi.TransactionContextInterface, fileID string) (string, error) {
	if fileID == "" {
		return "", fmt.Errorf("fileID must not be empty")
	}

	fileHash, err := ctx.GetStub().GetPrivateData("filehashCollection", fileID)
	if err != nil {
		return "", fmt.Errorf("failed to read file hash from private data collection: %w", err)
	}
	if fileHash == nil {
		return "", fmt.Errorf("file not found")
	}
	return string(fileHash), nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&FileHashContract{})
	if err != nil {
		log.Fatalf("Error creating chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Fatalf("Error starting chaincode: %v", err)
	}
}