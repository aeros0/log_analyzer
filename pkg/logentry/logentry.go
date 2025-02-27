package logentry

import (
	"fmt"
	"regexp"
	"time"
)

type LogEntry struct {
	Timestamp time.Time
	Level     string
	Message   string
}

func ParseLogLine(line string) (LogEntry, error) {
	//fmt.Println("lines:", line)
	re := regexp.MustCompile(`\[(.*?)\]\s+(.*?)\s+-\s+IP:.*?\s+(Error\s+500\s+-\s+.*)?`) // Corrected regex
	matches := re.FindStringSubmatch(line)

	//fmt.Println("Matches:", matches) // Debugging

	if len(matches) < 3 {
		return LogEntry{}, fmt.Errorf("invalid log format: %s", line)
	}

	timestamp, err := time.Parse(time.RFC3339, matches[1])
	if err != nil {
		return LogEntry{}, err
	}

	level := matches[2]
	message := ""
	if len(matches) > 3 && matches[3] != "" {
		message = matches[3]
	}

	return LogEntry{Timestamp: timestamp, Level: level, Message: message}, nil
}
