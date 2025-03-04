package stats

import (
	"log-analyzer/pkg/logentry"
	"math"
)

func GenerateStats(slidingWindow []logentry.LogEntry, patternCounts map[string]int, patternWeights map[string]float64, perSecondRates []int, skippedLogs int) map[string]interface{} {
	stats := make(map[string]interface{})
	var errorCount, infoCount, debugCount int
	for _, entry := range slidingWindow {
		switch entry.Level {
		case "ERROR":
			errorCount++
		case "INFO":
			infoCount++
		case "DEBUG":
			debugCount++
		}
	}

	total := len(slidingWindow)
	errorPercentage := float64(errorCount) / float64(total) * 100
	infoPercentage := float64(infoCount) / float64(total) * 100
	debugPercentage := float64(debugCount) / float64(total) * 100

	currentRate := CalculateRate(perSecondRates)
	peakRate := CalculatePeakRate(perSecondRates)

	// Declare errorRate here
	var errorRate float64

	if len(perSecondRates) > 0 {
		errorRate := float64(errorCount) / float64(len(perSecondRates))
		stats["errorRate"] = errorRate
	} else {
		stats["errorRate"] = 0.0
	}

	stats["entriesProcessed"] = total + skippedLogs
	stats["currentRate"] = currentRate
	stats["peakRate"] = peakRate
	stats["windowSize"] = len(slidingWindow)
	stats["errorPercentage"] = errorPercentage
	stats["infoPercentage"] = infoPercentage
	stats["debugPercentage"] = debugPercentage
	stats["errorCount"] = errorCount
	stats["infoCount"] = infoCount
	stats["debugCount"] = debugCount
	stats["errorRate"] = errorRate
	stats["skippedLogs"] = skippedLogs
	stats["patternCounts"] = patternCounts
	stats["patternWeights"] = patternWeights

	return stats
}

func CalculateRate(perSecondRates []int) float64 {
	sum := 0
	for _, rate := range perSecondRates {
		sum += rate
	}
	return float64(sum) / float64(len(perSecondRates))
}

func CalculatePeakRate(perSecondRates []int) float64 {
	peak := 0.0
	for _, rate := range perSecondRates {
		peak = math.Max(peak, float64(rate))
	}
	return peak
}
