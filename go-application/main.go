package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/status"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

func main() {
	// Step 1: Get the file path from the user
	filePath := "path/to/your/file.txt" // You can modify this part to read from command-line input or UI

	// Step 2: Open the file and hash it
	fileHash, err := hashFile(filePath)
	if err != nil {
		log.Fatalf("Failed to hash the file: %v", err)
	}

	// Step 3: Store the hash in CouchDB
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

// hashFile generates the SHA256 hash of a file
func hashFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Use SHA256 to hash the file content
	hash := sha256.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return "", fmt.Errorf("failed to hash file: %v", err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// storeHashInCouchDB stores the file hash in the CouchDB database
func storeHashInCouchDB(hash string) error {
	// Initialize CouchDB connection here (make sure CouchDB is running and accessible)
	// In production, you would use an actual CouchDB client, like github.com/go-kivik/kivik

	// For demonstration, assume the connection and storage are successful.
	// Simulating the storage step
	fmt.Printf("Storing hash %s in CouchDB...\n", hash)
	// You can connect to your CouchDB instance and store the hash.
	return nil
}

// invokeChaincode invokes the chaincode to store the hash and timestamp on the Fabric network
func invokeChaincode(hash string) error {
	// Set up the Fabric SDK
	sdk, err := fabsdk.New(fabsdk.WithConfig("connection-profile.yaml")) // Path to your Fabric connection profile
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

	// Prepare chaincode arguments
	args := []string{"storeImageHash", hash}

	// Invoke the chaincode
	response, err := client.Invoke(channel.Request{
		ChaincodeID: "mycc",       // Your chaincode name
		Fcn:           "storeImageHash",
		Args:          args,
	})
	if err != nil {
		return fmt.Errorf("failed to invoke chaincode: %v", err)
	}

	if response.Status != 200 {
		return fmt.Errorf("chaincode invocation failed with status: %d", response.Status)
	}

	fmt.Printf("Successfully invoked chaincode: %s\n", response.Payload)
	return nil
}
