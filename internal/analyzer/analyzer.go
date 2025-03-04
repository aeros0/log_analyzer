package analyzer

import (
	"fmt"
	"log-analyzer/pkg/logentry"
	"sync"
	"time"
)

var patternCountsMutex sync.RWMutex
var patternWeightsMutex sync.RWMutex

func ProcessLogs(logChan <-chan logentry.LogEntry, statsChan chan<- map[string]interface{}) {
	var lastLogTime time.Time
	var currentRate int
	var peakRate int
	var totalEntries int
	var errorCount int
	var errorCountperSec int
	var infoCount int
	var debugCount int
	var isFirst bool = true
	var errorCounts = make(map[string]int)
	var slidingWindow []logentry.LogEntry
	var adaptiveWindow int
	patternCounts := make(map[string]int)
	patternWeights := make(map[string]float64)
	// Calculate error percentages
	errorPercentages := make(map[string]float64)

	ticker := time.NewTicker(time.Second) // System-time ticker
	defer ticker.Stop()

	for {
		select {
		case entry := <-logChan:
			totalEntries++
			if isFirst {
				lastLogTime = entry.Timestamp
				isFirst = !isFirst
				fmt.Print(lastLogTime)
			}
			slidingWindow = append(slidingWindow, entry)
			if entry.Level == "ERROR" {
				//patternCountsMutex.Lock()
				patternCounts[entry.Message]++
				//patternCountsMutex.Unlock()
				//patternWeightsMutex.Lock()
				updateWeights(patternWeights, patternCounts, entry.Message)
				//patternWeightsMutex.Unlock()
			}

			if entry.Timestamp.Sub(lastLogTime) < time.Second {

				for msg, count := range errorCounts {
					errorPercentages[msg] = float64(count) / float64(totalEntries) * 100
				}

				//lastLogTime = entry.Timestamp

				errorCounts = make(map[string]int)
				slidingWindow = []logentry.LogEntry{}
				// } else if lastLogTime.IsZero() {
			} else {
				if currentRate > peakRate {
					peakRate = currentRate
				}
				lastLogTime = entry.Timestamp
				fmt.Print(lastLogTime)
				currentRate = 0
				errorCountperSec = 0
			}

			currentRate++
			if entry.Level == "ERROR" {
				errorCounts[entry.Message]++
				errorCount++
				errorCountperSec++
			}
			if entry.Level == "INFO" {
				infoCount++
			}
			if entry.Level == "DEBUG" {
				debugCount++
			}
		case <-ticker.C: // System-time trigger
			// Snapshot current statistics
			if currentRate > peakRate {
				peakRate = currentRate
			}

			// Calculate error percentages
			// errorPercentages := make(map[string]float64)
			// for msg, count := range errorCounts {
			// 	errorPercentages[msg] = float64(count) / float64(totalEntries) * 100
			// }

			statsData := make(map[string]interface{})
			// Generate and send stats
			//statsData := stats.GenerateStats(slidingWindow, patternCounts, patternWeights, nil, 0) // perSecondRates are irrelevant now
			statsData["currentRate"] = float64(currentRate)
			statsData["peakRate"] = float64(peakRate)
			statsData["totalEntries"] = float64(totalEntries)
			statsData["errorPercentages"] = errorPercentages
			statsData["entriesProcessed"] = int(totalEntries)
			statsData["windowSize"] = float64(adaptiveWindow)
			statsData["errorPercentage"] = float64(errorCount) / float64(totalEntries) * 100
			statsData["infoPercentage"] = float64(infoCount) / float64(totalEntries) * 100
			statsData["debugPercentage"] = float64(debugCount) / float64(totalEntries) * 100
			statsData["errorCount"] = errorCount
			statsData["infoCount"] = infoCount
			statsData["debugCount"] = debugCount
			statsData["errorRate"] = float64(errorCountperSec) / float64(currentRate) * 100
			//statsData["skippedLogs"] = skippedLogs
			statsData["patternCounts"] = patternCounts
			statsData["patternWeights"] = patternWeights

			//patternWeightsMutex.Lock()
			statsChan <- statsData
			//patternWeightsMutex.Unlock()
		}
	}
}

func updateWeights(patternWeights map[string]float64, patternCounts map[string]int, message string) {
	patternWeightsMutex.Lock()
	defer patternWeightsMutex.Unlock()

	if _, ok := patternWeights[message]; !ok {
		patternWeights[message] = 1.0
	}

	rateOfChange := calculateRateOfChange(patternCounts, message, 10*time.Second)

	if rateOfChange >= 4.0 {
		patternWeights[message] *= 3.0
	}
}

func calculateRateOfChange(patternCounts map[string]int, message string, window time.Duration) float64 {
	patternCountsMutex.RLock()
	defer patternCountsMutex.RUnlock()

	countNow := patternCounts[message]
	countBefore := 0 // Implement logic to get count from 'window' time ago
	if countBefore == 0 {
		return 0.0
	}

	return float64(countNow) / float64(countBefore)
}
