package hashlib

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

// Default number of threads
const NUM_THREADS = 4

// GetSHA computes SHA256 hash for a chunk of a file
func GetSHA(filename string, startIndex int64, bytesPerThread int64, wg *sync.WaitGroup, hashResults chan<- string) {
	defer wg.Done()

	absPath, err := filepath.Abs(filename)
	if err != nil {
		log.Println("Error resolving file path:", err)
		return
	}

	file, err := os.Open(absPath)
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	_, err = file.Seek(startIndex, 0)
	if err != nil {
		log.Println("Error seeking file:", err)
		return
	}

	hasher := sha256.New()
	buffer := make([]byte, 4096)
	var totalRead int64

	for totalRead < bytesPerThread {
		n, err := file.Read(buffer)
		if n > 0 {
			hasher.Write(buffer[:n])
			totalRead += int64(n)
		}
		if err == io.EOF || totalRead >= bytesPerThread {
			break
		} else if err != nil {
			log.Println("Error reading file:", err)
			return
		}
	}

	hashResults <- hex.EncodeToString(hasher.Sum(nil))
}

// HashFile generates a SHA256 hash of an entire file
func HashFile(filePath string) (string, error) {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve file path: %s", filePath)
	}

	file, err := os.Open(absPath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Get file size
	fileInfo, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to get file size: %v", err)
	}
	fileSize := fileInfo.Size()

	// Split work into threads
	bytesPerThread := fileSize / int64(NUM_THREADS)
	var wg sync.WaitGroup
	hashResults := make(chan string, NUM_THREADS)

	for i := 0; i < NUM_THREADS; i++ {
		wg.Add(1)
		startIndex := int64(i) * bytesPerThread
		go GetSHA(filePath, startIndex, bytesPerThread, &wg, hashResults)
	}

	wg.Wait()
	close(hashResults)

	// Collect results and ensure deterministic hashing order
	var hashes []string
	for hash := range hashResults {
		hashes = append(hashes, hash)
	}
	sort.Strings(hashes) // Sort to ensure consistent order

	// Combine hashes
	finalHasher := sha256.New()
	for _, hash := range hashes {
		finalHasher.Write([]byte(hash))
	}

	return hex.EncodeToString(finalHasher.Sum(nil)), nil
}