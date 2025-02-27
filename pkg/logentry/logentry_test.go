package logentry

import (
	"testing"
	"time"
)

func TestParseLogLine(t *testing.T) {
	testCases := []struct {
		input    string
		expected LogEntry
		wantErr  bool
	}{
		{
			input: "[2025-02-27T12:01:44Z] ERROR - IP:192.168.146.208 Error 500 - Null pointer exception",
			expected: LogEntry{
				Timestamp: time.Date(2025, 2, 27, 12, 1, 44, 0, time.UTC),
				Level:     "ERROR",
				Message:   "Error 500 - Null pointer exception",
			},
			wantErr: false,
		},
		{
			input: "[2025-02-27T12:01:43Z] INFO - IP:192.168.8.205",
			expected: LogEntry{
				Timestamp: time.Date(2025, 2, 27, 12, 1, 43, 0, time.UTC),
				Level:     "INFO",
				Message:   "",
			},
			wantErr: false,
		},
		{
			input:   "invalid log line",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		result, err := ParseLogLine(tc.input)
		if (err != nil) != tc.wantErr {
			t.Errorf("ParseLogLine(%q) error = %v, wantErr %v", tc.input, err, tc.wantErr)
			continue
		}
		if !tc.wantErr && result != tc.expected {
			t.Errorf("ParseLogLine(%q) = %+v, want %+v", tc.input, result, tc.expected)
		}
	}
}
