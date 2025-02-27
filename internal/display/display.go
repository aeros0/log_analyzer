package display

import (
	"fmt"
	"os"
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
			//fmt.Println(stats) // Debugging: Print stats
			displayReport(stats)
		case <-ticker.C:
			// If needed, display last stats again
			//fmt.Println("Ticker tick")
		case <-done:
			return
		}
	}
}

func displayReport(stats map[string]interface{}) {
	// Implement formatting and printing of the report
	fmt.Printf("Log Analysis Report (Last Updated: %s)\n", time.Now().UTC().Format(time.RFC3339))
	os.Stdout.Sync() // Flush the output
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	os.Stdout.Sync()
	fmt.Printf("Runtime Stats:\n")
	os.Stdout.Sync()
	fmt.Printf("• Entries Processed: %d\n", stats["entriesProcessed"])
	os.Stdout.Sync()
	fmt.Printf("• Current Rate: %.2f entries/sec (Peak: %.2f entries/sec)\n", stats["currentRate"], stats["peakRate"])
	os.Stdout.Sync()
	fmt.Printf("• Adaptive Window: %d sec\n", stats["windowSize"])
	os.Stdout.Sync()
	fmt.Println("\nPattern Analysis:")
	os.Stdout.Sync()
	fmt.Printf("• ERROR: %.2f%% (%d entries)\n", stats["errorPercentage"], stats["errorCount"])
	os.Stdout.Sync()
	fmt.Printf("• INFO: %.2f%% (%d entries)\n", stats["infoPercentage"], stats["infoCount"])
	os.Stdout.Sync()
	fmt.Printf("• DEBUG: %.2f%% (%d entries)\n", stats["debugPercentage"], stats["debugCount"])
	os.Stdout.Sync()
	fmt.Println("\nDynamic Insights:")
	os.Stdout.Sync()
	fmt.Printf("• Error Rate: %.2f errors/sec\n", stats["errorRate"])
	os.Stdout.Sync()
	// Add other dynamic insights and self-evolving alerts
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	os.Stdout.Sync()
	fmt.Println("Press Ctrl+C to exit")
	os.Stdout.Sync()
}
