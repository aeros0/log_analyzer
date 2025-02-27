package stats

import (
	"testing"
)

func TestCalculateRate(t *testing.T) {
	rates := []int{10, 20, 30}
	expected := 20.0
	result := CalculateRate(rates)
	if result != expected {
		t.Errorf("CalculateRate(%v) = %f, want %f", rates, result, expected)
	}
}

func TestCalculatePeakRate(t *testing.T) {
	rates := []int{10, 50, 20, 40}
	expected := 50.0
	result := CalculatePeakRate(rates)
	if result != expected {
		t.Errorf("CalculatePeakRate(%v) = %f, want %f", rates, result, expected)
	}
}

func TestGenerateStats(t *testing.T) {
	// Add test cases for GenerateStats as needed
	// ...
}

func TestAdjustWindow(t *testing.T) {
	// Add test cases for AdjustWindow as needed
	// ...
}
