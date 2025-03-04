package display

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

var patternCountsMutex sync.RWMutex
var patternWeightsMutex sync.RWMutex

func DisplayStats(statsChan <-chan map[string]interface{}, done chan bool) {
	for {
		select {
		case stats, ok := <-statsChan:
			if !ok {
				return
			}
			clearScreen()
			displayReport(stats)
		case <-done:
			return
		}
	}
}

// ...

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func displayReport(stats map[string]interface{}) {
	now := time.Now().UTC()
	fmt.Printf("Log Analysis Report (Last Updated: %s)\n", now.Format("2006-01-02 15:04:05 UTC"))
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Runtime Stats:\n")
	fmt.Printf("• Entries Processed: %d\n", stats["entriesProcessed"])
	fmt.Printf("• Current Rate: %.0f entries/sec (Peak: %.0f entries/sec)\n", stats["currentRate"], stats["peakRate"])
	fmt.Printf("• Adaptive Window: %d sec\n", stats["windowSize"])

	fmt.Println("\nPattern Analysis:")
	fmt.Printf("• ERROR: %.0f%% (%d entries)\n", stats["errorPercentage"], stats["errorCount"])
	fmt.Printf("• INFO: %.0f%% (%d entries)\n", stats["infoPercentage"], stats["infoCount"])
	fmt.Printf("• DEBUG: %.0f%% (%d entries)\n", stats["debugPercentage"], stats["debugCount"])

	fmt.Println("\nDynamic Insights:")
	fmt.Printf("• Error Rate: %.1f errors/sec\n", stats["errorRate"])

	// Emerging Pattern (Placeholder)
	fmt.Println("• Emerging Pattern: (Placeholder)")

	// Top Errors (Placeholder)
	fmt.Println("• Top Errors:")
	patternCountsCopy := make(map[string]int)
	patternWeightsCopy := make(map[string]float64)

	patternCountsMutex.RLock()
	originalPatternCounts, ok := stats["patternCounts"].(map[string]int)
	if ok {
		for k, v := range originalPatternCounts {
			patternCountsCopy[k] = v
		}
	}
	patternCountsMutex.RUnlock()

	patternWeightsMutex.RLock()
	originalPatternWeights, ok := stats["patternWeights"].(map[string]float64)
	if ok {
		for k, v := range originalPatternWeights {
			patternWeightsCopy[k] = v
		}
	}
	patternWeightsMutex.RUnlock()

	topErrors := getTopErrors(patternCountsCopy, patternWeightsCopy, 3)
	for i, errorMsg := range topErrors {
		fmt.Printf("  %d. %s\n", i+1, errorMsg)
	}

	fmt.Println("\nSelf-Evolving Alerts:")
	// Self-Evolving Alerts (Placeholder)
	fmt.Println("  (Placeholder)")

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("Press Ctrl+C to exit")
}

func getTopErrors(patternCounts map[string]int, patternWeights map[string]float64, limit int) []string {
	type errorEntry struct {
		Message string
		Weight  float64
	}

	var errors []errorEntry
	patternCountsMutex.RLock()  // Read lock
	patternWeightsMutex.RLock() // Read lock
	for msg, count := range patternCounts {
		errors = append(errors, errorEntry{Message: msg, Weight: float64(count) * patternWeights[msg]})
	}
	patternWeightsMutex.RUnlock() // Read unlock
	patternCountsMutex.RUnlock()

	patternWeightsMutex.RLock()
	sort.Slice(errors, func(i, j int) bool {
		return errors[i].Weight > errors[j].Weight
	})
	patternWeightsMutex.RUnlock()

	var topErrors []string
	for i := 0; i < len(errors) && i < limit; i++ {
		topErrors = append(topErrors, errors[i].Message)
	}

	return topErrors
}
