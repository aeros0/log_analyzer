package main

import (
	"log-analyzer/internal/analyzer"
	"log-analyzer/internal/display"
	"log-analyzer/internal/reader"
	"log-analyzer/pkg/logentry"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	logChan := make(chan logentry.LogEntry, 10000)
	statsChan := make(chan map[string]interface{})
	done := make(chan bool)
	var wg sync.WaitGroup

	wg.Add(3) // Add 3 goroutines to wait for

	go func() {
		reader.ReadLogs(logChan) //Removed scanner from argument.
		wg.Done()
	}()

	go func() {
		analyzer.ProcessLogs(logChan, statsChan)
		wg.Done()
	}()

	go func() {
		display.DisplayStats(statsChan, done)
		wg.Done()
	}()

	// Handle Ctrl+C
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChan
		done <- true
	}()

	<-done
	wg.Wait() // Wait for all goroutines to finish
}
