package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Octavian-Anghel/Capstone-Project/hashlib"
)

func main() {
	// Ensure a file path is provided as an argument
	if len(os.Args) < 2 {
		log.Fatal("Usage: ./Capstone-Project <file_path>")
	}

	filePath := os.Args[1]

	// Compute hash
	hash, err := hashlib.HashFile(filePath)
	if err != nil {
		log.Fatalf("Error hashing file: %v", err)
	}

	fmt.Printf("SHA256 Hash of %s: %s\n", filePath, hash)
}