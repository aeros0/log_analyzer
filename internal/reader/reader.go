package reader

import (
	"bufio"
	"fmt"
	"log-analyzer/pkg/logentry"
	"os"
	"time"
)

func ReadLogs(logChan chan<- logentry.LogEntry) {
	fmt.Println("Log reader started")

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	var lastModTime time.Time
	var lastOffset int64

	for range ticker.C {
		fileInfo, err := os.Stat("test_logs.log")
		if err != nil {
			fmt.Println("Error getting file info:", err)
			continue
		}

		if fileInfo.ModTime().After(lastModTime) {
			file, err := os.Open("test_logs.log")
			if err != nil {
				fmt.Println("Error opening log file:", err)
				continue
			}

			_, err = file.Seek(lastOffset, 0)
			if err != nil {
				fmt.Println("Error seeking file:", err)
				file.Close()
				continue
			}

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				entry, err := logentry.ParseLogLine(line)
				if err != nil {
					fmt.Println("Error reading line:", err)
					continue
				}
				logChan <- entry
				lastOffset, _ = file.Seek(0, 1)
			}

			if err := scanner.Err(); err != nil {
				fmt.Println("Error scanning log file:", err)
			}
			lastModTime = fileInfo.ModTime()
			file.Close()
		}
	}

	fmt.Println("Log channel closed")
	fmt.Println("Log reader finished")
	close(logChan)
}
