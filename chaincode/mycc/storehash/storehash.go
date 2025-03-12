package main

import (
	"fmt"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

// SimpleChaincode structure
type SimpleChaincode struct{}

// Init initializes the chaincode
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

// Invoke function is called to invoke chaincode operations
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()

	if function == "storeImageHash" {
		return t.storeImageHash(stub, args)
	} else if function == "queryImageHash" {
		return t.queryImageHash(stub, args)
	}

	return shim.Error("Invalid function name")
}

// storeImageHash stores an image hash with a timestamp
func (t *SimpleChaincode) storeImageHash(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	imageHash := args[0]
	timestamp := time.Now().UTC().Format(time.RFC3339)

	key := fmt.Sprintf("image_%s", imageHash)
	imageData := fmt.Sprintf(`{"imageHash":"%s", "timestamp":"%s"}`, imageHash, timestamp)

	err := stub.PutState(key, []byte(imageData))
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to store image hash: %s", err.Error()))
	}

	return shim.Success([]byte(fmt.Sprintf("Stored successfully: %s", key)))
}

// queryImageHash retrieves an image hash
func (t *SimpleChaincode) queryImageHash(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
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

// main function to start chaincode
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting SimpleChaincode: %s", err)
	}
}