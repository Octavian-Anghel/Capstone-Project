package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

// SimpleChaincode structure for the chaincode
type SimpleChaincode struct{}

// Init initializes the chaincode
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

// Invoke function is called to invoke chaincode operations
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()

	switch function {
	case "storeImageHash":
		return t.storeImageHash(stub, args)
	case "queryImageHash":
		return t.queryImageHash(stub, args)
	default:
		return shim.Error("Invalid function name")
	}
}

// storeImageHash stores an image hash along with the current timestamp
func (t *SimpleChaincode) storeImageHash(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 || args[0] == "" {
		return shim.Error("Incorrect number of arguments. Expecting 1: Image Hash")
	}

	imageHash := args[0]
	timestamp := time.Now().UTC().Format(time.RFC3339)

	data := map[string]string{
		"imageHash": imageHash,
		"timestamp": timestamp,
	}

	imageData, err := json.Marshal(data)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to serialize data: %s", err.Error()))
	}

	key := fmt.Sprintf("image_%s", imageHash)
	if err := stub.PutState(key, imageData); err != nil {
		return shim.Error(fmt.Sprintf("Failed to store image hash: %s", err.Error()))
	}

	return shim.Success([]byte(fmt.Sprintf("Image hash and timestamp stored successfully: %s", key)))
}

// queryImageHash retrieves the image hash and timestamp
func (t *SimpleChaincode) queryImageHash(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 || args[0] == "" {
		return shim.Error("Incorrect number of arguments. Expecting 1: Image Hash")
	}

	imageHash := args[0]
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
		log.Fatalf("Error starting SimpleChaincode: %v", err)
	}
}