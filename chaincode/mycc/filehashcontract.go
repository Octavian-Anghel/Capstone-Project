package main

import (
    "fmt"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// FileHashContract defines the smart contract
type FileHashContract struct {
    contractapi.Contract
}

// StoreFileHash stores the hash of a file in the private data collection
func (s *FileHashContract) StoreFileHash(ctx contractapi.TransactionContextInterface, fileID string, fileHash string) error {
    // Store the file hash in the private data collection
    err := ctx.GetStub().PutPrivateData("filehashCollection", fileID, []byte(fileHash))
    if err != nil {
        return fmt.Errorf("failed to store file hash in private data collection: %v", err)
    }
    return nil
}

// QueryFileHash retrieves a file hash from the private data collection
func (s *FileHashContract) QueryFileHash(ctx contractapi.TransactionContextInterface, fileID string) (string, error) {
    // Retrieve the file hash from the private data collection
    fileHash, err := ctx.GetStub().GetPrivateData("filehashCollection", fileID)
    if err != nil {
        return "", fmt.Errorf("failed to read file hash from private data collection: %v", err)
    }
    if fileHash == nil {
        return "", fmt.Errorf("file not found")
    }
    return string(fileHash), nil
}

func main() {
    chaincode, err := contractapi.NewChaincode(&FileHashContract{})
    if err != nil {
        fmt.Printf("Error creating chaincode: %s", err.Error())
        return
    }

    if err := chaincode.Start(); err != nil {
        fmt.Printf("Error starting chaincode: %s", err.Error())
    }
}
