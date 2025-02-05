package main

import (
	"fmt"
	"time"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

// SimpleChaincode structure for the chaincode
type SimpleChaincode struct {
}

// Init initializes the chaincode
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

// Invoke function is called to invoke chaincode operations
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()

	if function == "storeImageHash" {
		// Store the image hash with the timestamp
		return t.storeImageHash(stub, args)
	}

	// If the function is not found, return an error
	return shim.Error("Invalid function name")
}

// storeImageHash stores an image hash along with the current timestamp
func (t *SimpleChaincode) storeImageHash(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// Validate input args
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1: Image Hash")
	}

	imageHash := args[0]
	timestamp := time.Now().UTC().Format(time.RFC3339)

	// Generate a unique key for the image hash entry (can be imageHash or a new key)
	key := fmt.Sprintf("image_%s", imageHash)

	// Store the image hash and timestamp in a composite JSON format (or any format)
	imageData := fmt.Sprintf(`{"imageHash":"%s", "timestamp":"%s"}`, imageHash, timestamp)

	// Save the data to the ledger (CouchDB)
	err := stub.PutState(key, []byte(imageData))
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to store image hash: %s", err.Error()))
	}

	return shim.Success([]byte(fmt.Sprintf("Image hash and timestamp stored successfully: %s", key)))
}

// QueryImageHash retrieves the image hash and timestamp based on the image hash
func (t *SimpleChaincode) queryImageHash(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// Validate input args
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1: Image Hash")
	}

	imageHash := args[0]

	// Get the state from the ledger
	key := fmt.Sprintf("image_%s", imageHash)
	imageData, err := stub.GetState(key)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to get image data: %s", err.Error()))
	}

	if imageData == nil {
		return shim.Error("Image hash not found")
	}

	return shim.Success(imageData)
}

// main function to start the chaincode
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting SimpleChaincode: %s", err)
	}
}
