package analyzer

import (
	"fmt"
	"log-analyzer/internal/stats"
	"log-analyzer/pkg/logentry"
	"sync/atomic"
	"time"
)

func ProcessLogs(logChan <-chan logentry.LogEntry, statsChan chan<- map[string]interface{}) {
	//fmt.Println("Log processor started")
	slidingWindow := make([]logentry.LogEntry, 0, 60)
	patternCounts := make(map[string]int)
	patternWeights := make(map[string]float64)
	perSecondRates := make([]int, 60)
	lastUpdateTime := time.Now()
	var skippedLogs int32

	for entry := range logChan {
		//fmt.Println("Log entry received")
		if entry.Timestamp.IsZero() {
			atomic.AddInt32(&skippedLogs, 1)
			fmt.Println("Skipped logs:", atomic.LoadInt32(&skippedLogs))
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
			//fmt.Println("Last update:", lastUpdateTime, "Now:", now)
		}
		perSecondRates[len(perSecondRates)-1]++

		adjustWindow(slidingWindow, perSecondRates)

		timeDiff := time.Since(lastUpdateTime)
		//fmt.Println("Time Diff:", timeDiff, "Last Update:", lastUpdateTime, "Now:", now) //Debug
		if timeDiff < time.Second {
			stats := stats.GenerateStats(slidingWindow, patternCounts, patternWeights, perSecondRates, int(atomic.LoadInt32(&skippedLogs)))
			statsChan <- stats
			fmt.Println("Stats sent")
			lastUpdateTime = time.Now()
		}
	}
	close(statsChan)
}

func adjustWindow(slidingWindow []logentry.LogEntry, perSecondRates []int) {
	rate := stats.CalculateRate(perSecondRates)
	if rate > 2500 && len(slidingWindow) > 30 {
		slidingWindow = slidingWindow[len(slidingWindow)-30:]
	} else if rate < 600 && len(slidingWindow) < 120 {
		//No implementation for increase window size.
	}
}

func updateWeights(patternWeights map[string]float64, patternCounts map[string]int, message string) {
	// Implement logic for updating weights based on rate of change
	// ...
}
