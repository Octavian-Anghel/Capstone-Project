package hashlib

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"sync"
	"testing"
)

// Helper function to compute SHA-256 hash
func computeSHA256(data []byte) string {
	hash := sha256.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}

// Test Full-File Hashing
func TestHashFile(t *testing.T) {
	testFile := "testfile.txt"
	content := []byte("Hello, Hyperledger Fabric Parallel Hashing!")
	expectedHash := computeSHA256(content)

	err := os.WriteFile(testFile, content, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(testFile)

	hash, err := HashFile(testFile)
	if err != nil {
		t.Fatalf("Error hashing file: %v", err)
	}

	if hash != expectedHash {
		t.Errorf("Expected hash %s, but got %s", expectedHash, hash)
	}
}

// Test Handling of Missing Files
func TestHashFile_FileNotFound(t *testing.T) {
	_, err := HashFile("nonexistent.txt")
	if err == nil {
		t.Errorf("Expected error for missing file, but got nil")
	}
}

// Test Parallel Hashing with Multiple Threads
func TestParallelHashing(t *testing.T) {
	testFile := "parallel_testfile.txt"
	content := []byte("Parallel Hashing Test in Go!")
	expectedHash := computeSHA256(content)

	err := os.WriteFile(testFile, content, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(testFile)

	var wg sync.WaitGroup
	hashResults := make(chan string, 4)
	fileSize := int64(len(content))
	bytesPerThread := fileSize / 4 // 4 threads

	for i := 0; i < 4; i++ {
		startIndex := int64(i) * bytesPerThread
		wg.Add(1)
		go GetSHA(testFile, startIndex, bytesPerThread, &wg, hashResults)
	}

	wg.Wait()
	close(hashResults)

	var hashes []string
	for hash := range hashResults {
		hashes = append(hashes, hash)
	}

	// Ensure all chunks produced valid hashes
	if len(hashes) != 4 {
		t.Errorf("Expected 4 hashes, got %d", len(hashes))
	}

	// Hashing the entire file should still match
	fullHash, err := HashFile(testFile)
	if err != nil {
		t.Fatalf("Error hashing entire file: %v", err)
	}

	if fullHash != expectedHash {
		t.Errorf("Expected full hash %s, but got %s", expectedHash, fullHash)
	}
}