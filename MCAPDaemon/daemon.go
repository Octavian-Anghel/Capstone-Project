package main

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strings"
	"sync"
	"time"
	"github.com/fsnotify/fsnotify"
)

func printTime(format string, a ...interface{}) {
		fmt.Printf("[%s] ", time.Now().Format(time.RFC3339))
	fmt.Printf(format+"\n", a...)
}

func exit(format string, a ...interface{}) {
		printTime(format, a...)
	os.Exit(1)
}

func getMagicBytes(filePath string) (string, error) {
		file, err := os.Open(filePath)
	if err != nil {
			return "", err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
			return "", err
	}
	fileSize := fileInfo.Size()

	if fileSize < 5 {
			return "", errors.New("filesize too small, likely bad file")
	}

	start := fileSize - int64(7)
	buff := make([]byte, 7)

	_, err = file.ReadAt(buff, start)
	if err != nil {
			return "", err
	}

	return string(buff), nil
}

func dedupLoop(w *fsnotify.Watcher) {
		var (
		waitFor    = 100 * time.Millisecond
		mu         sync.Mutex
		timers     = make(map[string]*time.Timer)
		printEvent = func(e fsnotify.Event) {
				printTime("Detected event: %s", e)

			if strings.HasSuffix(e.Name, ".mcap") {
					magic, err := getMagicBytes(e.Name)
				if err != nil {
						fmt.Printf("Failed to read magic bytes from %s: %v\n", e.Name, err)
				} else {
						fmt.Printf("Magic bytes from %s: %s\n", e.Name, magic)
					if magic == "MCAP0\r\n" {
							fmt.Println("Valid MCAP file detected! pushing over to the hash and upload daemon\n")
						// Additional processing here
							HashNUpload(e.Name)
					}
				}
			}

			mu.Lock()
			delete(timers, e.Name)
			mu.Unlock()
		}
	)

	for {
			select {
			case err, ok := <-w.Errors:
			if !ok {
					return
			}
			printTime("ERROR: %s", err)

		case e, ok := <-w.Events:
			if !ok {
					return
			}

			if !e.Has(fsnotify.Create) && !e.Has(fsnotify.Write) {
					continue
			}

			mu.Lock()
			t, ok := timers[e.Name]
			mu.Unlock()

			if !ok {
					t = time.AfterFunc(math.MaxInt64, func() { printEvent(e) })
				t.Stop()
				mu.Lock()
				timers[e.Name] = t
				mu.Unlock()
			}

			t.Reset(waitFor)
		}
	}
}

func main() {
		w, err := fsnotify.NewWatcher()
	if err != nil {
			exit("creating a new watcher: %s", err)
	}
	defer w.Close()

	go dedupLoop(w)

	path := "/shared"
	err = w.Add(path)
	if err != nil {
			exit("%q: %s", path, err)
	}

	// Prevent main from exiting
	select {}
}
