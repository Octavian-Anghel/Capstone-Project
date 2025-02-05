package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"sync"
)

const NUM_THREADS = 4

func getSHA(filename string, startIndex int64, bytesPerThread int64, wg *sync.WaitGroup, hashResults chan<- string) {
	defer wg.Done()

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Move to the correct position
	_, err = file.Seek(startIndex, 0)
	if err != nil {
		fmt.Println("Error seeking file:", err)
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
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}
	}

	hashResults <- fmt.Sprintf("%x", hasher.Sum(nil))
	fmt.Println("Hashing complete for chunk starting at:", startIndex)
}

func main() {
	filename := "printdata.mcap"
	fileInfo, err := os.Stat(filename)
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return
	}

	fileSize := fileInfo.Size()
	bytesPerThread := fileSize / NUM_THREADS

	var wg sync.WaitGroup
	hashResults := make(chan string, NUM_THREADS)

	for i := 0; i < NUM_THREADS; i++ {
		startIndex := int64(i) * bytesPerThread
		wg.Add(1)
		go getSHA(filename, startIndex, bytesPerThread, &wg, hashResults)
	}

	wg.Wait()
	close(hashResults)

	var hashes []string
	for hash := range hashResults {
		hashes = append(hashes, hash)
	}

	fmt.Println("Hashes:", hashes)
}
