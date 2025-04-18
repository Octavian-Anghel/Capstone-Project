package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"sync"
	"sort"
)

const NUM_THREADS = 56

type HashResult struct {
	Index int
	Hash string
}

func GetSHA(filename string, startIndex int64, bytesPerThread int64, index int, wg *sync.WaitGroup, hashResults chan<- HashResult) {
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

	hash :=fmt.Sprintf("%x", hasher.Sum(nil))
	hashResults <- HashResult{Index: index, Hash: hash}
	//fmt.Println("Hashing complete for chunk starting at:", startIndex)
}

func HashNUpload(fn string) string {
	filename := fn
	fileInfo, err := os.Stat(filename)
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return ""
	}

	fileSize := fileInfo.Size()
	bytesPerThread := fileSize / NUM_THREADS

	var wg sync.WaitGroup
	hashResults := make(chan HashResult, NUM_THREADS)

	for i := 0; i < NUM_THREADS; i++ {
		startIndex := int64(i) * bytesPerThread
		wg.Add(1)
		go GetSHA(filename, startIndex, bytesPerThread, i, &wg, hashResults)
	}

	wg.Wait()
	close(hashResults)

	var results []HashResult

	for hr := range hashResults {
		results = append(results, hr)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Index < results[j].Index
	})

	finalHash := ""
	for _, r := range results {
		finalHash += r.Hash
	}

	return finalHash
}
