package reader

import (
	"bufio"
	"fmt"
	"log-analyzer/pkg/logentry"
	"os"
	"time"
)

func ReadLogs(logChan chan<- logentry.LogEntry, statsChan chan<- map[string]interface{}, done chan bool) {
	fmt.Println("Log reader started")

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	var lastModTime time.Time
	var lastOffset int64

	for {
		select {
		case <-ticker.C:
			startTime := time.Now() // Record start time

			fileInfo, err := os.Stat("test_logs.log")
			if err != nil {
				// Skip updates if file access fails
				continue
			}

			if fileInfo.ModTime().After(lastModTime) {
				file, err := os.Open("test_logs.log")
				if err != nil {
					// Skip updates if file open fails
					continue
				}

				_, err = file.Seek(lastOffset, 0)
				if err != nil {
					file.Close()
					continue
				}

				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					if time.Since(startTime) >= 900*time.Millisecond { // Stop reading 100ms before next tick
						break
					}
					line := scanner.Text()
					entry, err := logentry.ParseLogLine(line)
					if err != nil {
						continue
					}
					logChan <- entry
					lastOffset, _ = file.Seek(0, 1)
				}

				if err := scanner.Err(); err != nil {
					// Skip updates on scanning error
				}
				lastModTime = fileInfo.ModTime()
				file.Close()
			}
		case <-done:
			fmt.Println("Log channel closed")
			fmt.Println("Log reader finished")
			close(logChan)
			close(statsChan)
			return
		}
	}
}
