package hashlib

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

// Number of concurrent threads for hashing
const NUM_THREADS = 4

// GetSHA computes a SHA256 hash for a chunk of a file
func GetSHA(filename string, startIndex int64, bytesPerThread int64, wg *sync.WaitGroup, hashResults chan<- string) {
	defer wg.Done()

	file, err := os.Open(filename)
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
	log.Printf("Hashing complete for chunk starting at: %d\n", startIndex)
}

// HashFile generates the SHA256 hash of an entire file
func HashFile(filePath string) (string, error) {
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