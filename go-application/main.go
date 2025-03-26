package main

import (
	"fmt"
	"log"

	"github.com/Octavian-Anghel/Capstone-Project/hashlib"
)

func main() {
	filePath := "testfile.txt"

	// Hash full file
	hash, err := hashlib.HashFile(filePath)
	if err != nil {
		log.Fatalf("Failed to hash file: %v", err)
	}

	fmt.Printf("SHA-256 Hash of %s: %s\n", filePath, hash)
}