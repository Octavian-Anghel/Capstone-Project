package hashlib

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// Number of threads for parallel hashing
const NUM_THREADS = 4

// GetSHA computes SHA256 hash for a chunk of a file
func GetSHA(filename string, startIndex int64, bytesPerThread int64, wg *sync.WaitGroup, hashResults chan<- string) {
	defer wg.Done()

	// Resolve and clean the absolute path
	absPath, err := filepath.Abs(filename)
	if err != nil {
		log.Println("Failed to resolve absolute path:", err)
		return
	}
	cleanPath := filepath.Clean(absPath)

	// Ensure file exists and is not a directory
	fileInfo, err := os.Stat(cleanPath)
	if err != nil || fileInfo.IsDir() {
		log.Println("Invalid file path:", cleanPath)
		return
	}

	// Open the file securely
	file, err := os.Open(cleanPath) // Removed os.OpenFile for security
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	_, err = file.Seek(startIndex, io.SeekStart)
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
	log.Printf("Hashing complete for chunk starting at: %d\n", startIndex)
}

// HashFile generates SHA256 hash of an entire file
func HashFile(filePath string) (string, error) {
	// Resolve and clean the absolute path
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path: %w", err)
	}
	cleanPath := filepath.Clean(absPath)

	// Ensure file exists and is not a directory
	fileInfo, err := os.Stat(cleanPath)
	if err != nil || fileInfo.IsDir() {
		return "", fmt.Errorf("invalid file path: %s", cleanPath)
	}

	// Open file securely in read-only mode
	file, err := os.Open(cleanPath) // Removed os.OpenFile for security
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	hash := sha256.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return "", fmt.Errorf("failed to hash file: %w", err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}