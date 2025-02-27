package reader

import (
	"bufio"
	"log-analyzer/pkg/logentry"
	"strings"
	"testing"
)

func TestReadLogs(t *testing.T) {
	input := "[2025-02-27T12:01:43Z] INFO - IP:192.168.8.205\n[2025-02-27T12:01:44Z] ERROR - IP:192.168.146.208 Error 500 - Null pointer exception"
	scanner := bufio.NewScanner(strings.NewReader(input))
	logChan := make(chan logentry.LogEntry, 2)
	go ReadLogs(scanner, logChan)

	entry1 := <-logChan
	if entry1.Level != "INFO" {
		t.Errorf("Expected INFO, got %s", entry1.Level)
	}

	entry2 := <-logChan
	if entry2.Level != "ERROR" {
		t.Errorf("Expected ERROR, got %s", entry2.Level)
	}

	_, ok := <-logChan
	if ok {
		t.Errorf("Channel should be closed")
	}
}
