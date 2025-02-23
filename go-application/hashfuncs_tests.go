package main

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"sync"
	"testing"
)

// Helper function to manually compute SHA-256
func computeSHA256(data []byte) string {
	hash := sha256.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}

// Test: Hashing a known file
func TestGetSHA_ValidFile(t *testing.T) {
	// Create a temporary test file
	testFile := "testdata/testfile.txt"
	content := []byte("Hello, Hyperledger Fabric!")
	expectedHash := computeSHA256(content)

	err := os.WriteFile(testFile, content, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(testFile)

	// Setup for parallel hashing
	var wg sync.WaitGroup
	hashResults := make(chan string, 1)

	wg.Add(1)
	go getSHA(testFile, 0, int64(len(content)), &wg, hashResults)

	wg.Wait()
	close(hashResults)

	// Get result
	hash := <-hashResults
	if hash != expectedHash {
		t.Errorf("Expected hash %s, but got %s", expectedHash, hash)
	}
}

// Test: Handling a missing file
func TestGetSHA_FileNotFound(t *testing.T) {
	var wg sync.WaitGroup
	hashResults := make(chan string, 1)

	wg.Add(1)
	go getSHA("nonexistent.txt", 0, 1024, &wg, hashResults)

	wg.Wait()
	close(hashResults)

	select {
	case hash := <-hashResults:
		t.Errorf("Expected no hash, but got %s", hash)
	default:
		// Test passes if no hash is received
	}
}

// Test: Handling an empty file
func TestGetSHA_EmptyFile(t *testing.T) {
	emptyFile := "testdata/empty.txt"
	err := os.WriteFile(emptyFile, []byte(""), 0644)
	if err != nil {
		t.Fatalf("Failed to create empty file: %v", err)
	}
	defer os.Remove(emptyFile)

	var wg sync.WaitGroup
	hashResults := make(chan string, 1)

	wg.Add(1)
	go getSHA(emptyFile, 0, 0, &wg, hashResults)

	wg.Wait()
	close(hashResults)

	hash := <-hashResults
	expectedHash := computeSHA256([]byte(""))
	if hash != expectedHash {
		t.Errorf("Expected empty file hash %s, but got %s", expectedHash, hash)
	}
}

// Test: Multi-threaded hash consistency
func TestGetSHA_MultiThread(t *testing.T) {
	testFile := "testdata/multithread.txt"
	content := []byte("Multithreaded hashing test")
	expectedHash := computeSHA256(content)

	err := os.WriteFile(testFile, content, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(testFile)

	var wg sync.WaitGroup
	hashResults := make(chan string, 4)

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go getSHA(testFile, 0, int64(len(content)), &wg, hashResults)
	}

	wg.Wait()
	close(hashResults)

	for hash := range hashResults {
		if hash != expectedHash {
			t.Errorf("Expected hash %s, but got %s", expectedHash, hash)
		}
	}
}