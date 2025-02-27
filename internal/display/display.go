package display

import (
	"fmt"
	"sort"
	"time"
)

func DisplayStats(statsChan <-chan map[string]interface{}, done chan bool) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case stats, ok := <-statsChan:
			if !ok {
				return
			}
			clearScreen()
			displayReport(stats)
		case <-ticker.C:
			// No action on ticker, display report only when stats are received
		case <-done:
			return
		}
	}
}

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
	topErrors := getTopErrors(stats["patternCounts"].(map[string]int), stats["patternWeights"].(map[string]float64), 3)
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
	for msg, count := range patternCounts {
		errors = append(errors, errorEntry{Message: msg, Weight: float64(count) * patternWeights[msg]})
	}

	sort.Slice(errors, func(i, j int) bool {
		return errors[i].Weight > errors[j].Weight
	})

	var topErrors []string
	for i := 0; i < len(errors) && i < limit; i++ {
		topErrors = append(topErrors, errors[i].Message)
	}

	return topErrors
}
