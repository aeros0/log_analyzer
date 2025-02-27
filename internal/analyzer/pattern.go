package analyzer

import (
	"math"
	"time"
)

//var mutex sync.Mutex

// func updateWeights(patternWeights map[string]float64, patternCounts map[string]int, message string) {
// 	mutex.Lock()
// 	defer mutex.Unlock()

// 	if _, ok := patternWeights[message]; !ok {
// 		patternWeights[message] = 1.0
// 	}

// 	rateOfChange := calculateRateOfChange(patternCounts, message, 10*time.Second)

// 	if rateOfChange >= 4.0 {
// 		patternWeights[message] *= 3.0
// 	}
// }

func calculateRateOfChange(patternCounts map[string]int, message string, window time.Duration) float64 {
	countNow := patternCounts[message]
	countBefore := 0 // Implement logic to get count from 'window' time ago
	if countBefore == 0 {
		return 0.0
	}

	return float64(countNow) / float64(countBefore)
}

func calculateTopErrors(patternCounts map[string]int, patternWeights map[string]float64, limit int) []string {
	// Implement logic to calculate the top errors based on their weighted scores.
	// Placeholder implementation:
	return []string{}
}

func calculatePeakRateError(perSecondRates []int) float64 {
	peak := 0.0
	for _, rate := range perSecondRates {
		peak = math.Max(peak, float64(rate))
	}
	return peak
}
