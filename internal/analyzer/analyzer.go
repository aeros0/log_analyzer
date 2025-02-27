package analyzer

import (
	"fmt"
	"log-analyzer/internal/stats"
	"log-analyzer/pkg/logentry"
	"sync"
	"sync/atomic"
	"time"
)

var mutex sync.Mutex

func ProcessLogs(logChan <-chan logentry.LogEntry, statsChan chan<- map[string]interface{}) {
	slidingWindow := make([]logentry.LogEntry, 0, 60)
	patternCounts := make(map[string]int)
	patternWeights := make(map[string]float64)
	perSecondRates := make([]int, 60)
	lastUpdateTime := time.Now()
	var skippedLogs int32
	bufferSize := 10000
	var errorCount int32
	var lastErrorRate float64

	for entry := range logChan {
		if entry.Timestamp.IsZero() {
			atomic.AddInt32(&skippedLogs, 1)
			continue
		}
		slidingWindow = append(slidingWindow, entry)
		if len(slidingWindow) > int(time.Since(slidingWindow[0].Timestamp).Seconds()) {
			slidingWindow = slidingWindow[1:]
		}

		if entry.Level == "ERROR" && entry.Message != "" {
			patternCounts[entry.Message]++
			updateWeights(patternWeights, patternCounts, entry.Message)
		}

		now := time.Now()
		if now.Second() != lastUpdateTime.Second() {
			perSecondRates = append(perSecondRates[1:], 0)
			lastUpdateTime = now
		}
		perSecondRates[len(perSecondRates)-1]++

		adjustWindow(slidingWindow, perSecondRates)

		if len(logChan) > bufferSize*9/10 {
			bufferSize = bufferSize * 3 / 2
			fmt.Printf("[%s] ⚠️ Burst detected: Resized buffer to %d\n", time.Now().UTC().Format("15:04:05"), bufferSize)
		}

		if entry.Level == "ERROR" {
			atomic.AddInt32(&errorCount, 1)
		}

		// Generate and send stats on every entry
		stats := stats.GenerateStats(slidingWindow, patternCounts, patternWeights, perSecondRates, int(atomic.LoadInt32(&skippedLogs)))
		statsChan <- stats

		timeDiff := time.Since(lastUpdateTime)
		if timeDiff >= time.Second {
			lastUpdateTime = time.Now()
			errorRate := stats["errorRate"].(float64)
			if errorRate > lastErrorRate*2 && lastErrorRate != 0 {
				fmt.Printf("[%s] ⚠️ High error rate (%.1f errors/sec), increased pattern weight\n", time.Now().UTC().Format("15:04:05"), errorRate)
				// Implement logic to increase pattern weights here
			}
			lastErrorRate = errorRate
		}
	}
	close(statsChan)
}

func adjustWindow(slidingWindow []logentry.LogEntry, perSecondRates []int) {
	rate := stats.CalculateRate(perSecondRates)
	if rate > 2500 && len(slidingWindow) > 30 {
		slidingWindow = slidingWindow[len(slidingWindow)-30:]
		fmt.Printf("[%s] ⚠️ Adjusted window to %d sec due to rate surge\n", time.Now().UTC().Format("15:04:05"), len(slidingWindow))
	} else if rate < 600 && len(slidingWindow) < 120 {
		// Implement logic to increase window size if needed
		// slidingWindow = slidingWindow[len(slidingWindow)-120:]
	}
}

func updateWeights(patternWeights map[string]float64, patternCounts map[string]int, message string) {
	mutex.Lock()
	defer mutex.Unlock()

	if _, ok := patternWeights[message]; !ok {
		patternWeights[message] = 1.0
	}

	rateOfChange := calculateRateOfChange(patternCounts, message, 10*time.Second)

	if rateOfChange >= 4.0 {
		patternWeights[message] *= 3.0
	}
}
