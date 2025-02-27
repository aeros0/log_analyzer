package analyzer

import (
	"math"
	"time"
)

// updateWeights updates the weights of error patterns based on their rate of change.
func updateWeightsOfChange(patternWeights map[string]float64, patternCounts map[string]int, message string) {
	// Implement logic for updating weights based on rate of change
	// Example: Track rate of change over a short window (e.g., 10 seconds)
	// and increase weight if frequency quadruples.

	// Placeholder logic - replace with your actual implementation
	if _, ok := patternWeights[message]; !ok {
		patternWeights[message] = 1.0 // Initial weight
	}

	//Simple race condition example.
	patternWeights[message] = patternWeights[message] + 0.1

	// More advanced logic would involve:
	// 1. Keeping track of historical counts over a sliding window.
	// 2. Calculating the rate of change.
	// 3. Applying a multiplier to the weight if the rate of change exceeds a threshold.
}

// calculateRateOfChange calculates the rate of change of an error pattern's frequency.
func calculateRateOfChange(patternCounts map[string]int, message string, window time.Duration) float64 {
	// Implement logic to calculate the rate of change over the specified window.
	// This would involve tracking counts over time and comparing them.
	// Placeholder implementation:
	return 0.0
}

// calculateErrorRate calculates the error rate for a given time window.
func calculateErrorRate(patternCounts map[string]int, window time.Duration) float64 {
	// Implement logic to calculate the error rate over the specified window.
	// Placeholder implementation:
	return 0.0
}

// calculateTopErrors calculates the top errors based on their weighted scores.
func calculateTopErrors(patternCounts map[string]int, patternWeights map[string]float64, limit int) []string {
	// Implement logic to calculate the top errors based on their weighted scores.
	// Placeholder implementation:
	return []string{}
}

// calculateEmergingPattern detects emerging patterns based on their rate of change.
func calculateEmergingPattern(patternCounts map[string]int, patternWeights map[string]float64, window time.Duration) string {
	// Implement logic to detect emerging patterns based on their rate of change.
	// Placeholder implementation:
	return ""
}

// calculatePeakRateError calculates the peak rate of errors in the last second.
func calculatePeakRateError(perSecondRates []int) float64 {
	peak := 0.0
	for _, rate := range perSecondRates {
		peak = math.Max(peak, float64(rate))
	}
	return peak
}
