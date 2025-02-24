package main

import (
	"fmt"
	"log"

	"github.com/Octavian-Anghel/Capstone-Project/hashlib"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

func main() {
	// Step 1: Get the file path
	filePath := "path/to/your/file.txt"

	// Step 2: Hash the entire file
	fileHash, err := hashlib.HashFile(filePath)
	if err != nil {
		log.Fatalf("Failed to hash the file: %v", err)
	}

	// Step 3: Store the hash in CouchDB (Simulated)
	err = storeHashInCouchDB(fileHash)
	if err != nil {
		log.Fatalf("Failed to store the hash in CouchDB: %v", err)
	}

	// Step 4: Invoke chaincode to upload hash to the blockchain
	err = invokeChaincode(fileHash)
	if err != nil {
		log.Fatalf("Failed to invoke chaincode: %v", err)
	}

	fmt.Println("File processed and uploaded successfully.")
}

// storeHashInCouchDB stores the file hash in the CouchDB database
func storeHashInCouchDB(hash string) error {
	// Simulated database storage
	fmt.Printf("Storing hash %s in CouchDB...\n", hash)
	return nil
}

// invokeChaincode invokes the chaincode to store the hash and timestamp on the Fabric network
func invokeChaincode(hash string) error {
	// Set up the Fabric SDK with the correct config provider
	configProvider := config.FromFile("connection-profile.yaml")
	sdk, err := fabsdk.New(configProvider)
	if err != nil {
		return fmt.Errorf("failed to create Fabric SDK: %v", err)
	}
	defer sdk.Close()

	// Create a client context
	clientContext := sdk.ChannelContext("mychannel", fabsdk.WithUser("Admin"), fabsdk.WithOrg("Org1"))

	// Create a channel client
	client, err := channel.New(clientContext)
	if err != nil {
		return fmt.Errorf("failed to create channel client: %v", err)
	}

	// Convert args to [][]byte (Fabric SDK requires byte slices)
	args := []string{"storeImageHash", hash}
	byteArgs := make([][]byte, len(args))
	for i, arg := range args {
		byteArgs[i] = []byte(arg)
	}

	// Invoke the chaincode
	response, err := client.Execute(channel.Request{
		ChaincodeID: "mycc",
		Fcn:         "storeImageHash",
		Args:        byteArgs,
	})
	if err != nil {
		return fmt.Errorf("failed to invoke chaincode: %v", err)
	}

	if response.TxValidationCode != 0 {
		return fmt.Errorf("chaincode invocation failed with validation code: %d", response.TxValidationCode)
	}

	fmt.Printf("Successfully invoked chaincode: %s\n", response.Payload)
	return nil
}