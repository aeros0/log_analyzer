# Log Analyzer

This tool analyzes log streams in real-time from standard input, providing runtime statistics, pattern analysis, and dynamic insights.

## Build and Run

1.  **Clone the repository:**

    ```bash
    git clone <repository_url>
    cd log-analyzer
    ```

2.  **Build the application:**

    ```bash
    go build ./cmd/log-analyzer
    ```

3.  **Run the application:**

    ```bash
    ./log-analyzer < test_log.log
    ```

    or

    ```bash
    ./log_generator.sh | ./cmd/log-analyzer
    ```

4.  **Exit the application:**

    Press `Ctrl+C`.

## Functionality

* **Self-Adjusting Time Window:** The tool dynamically adjusts the sliding window based on the processing rate.
* **Runtime Pattern Detection and Weighting:** Error patterns are tracked and weighted based on their frequency and rate of change.
* **Burst Handling and Resilience:** The tool handles log bursts without crashing.
* **Real-Time Display:** The console is updated every second with a report of the analysis.
* **Error Handling:** Malformed logs are skipped and counted.

## Concurrency Bug

A race condition is present in the `updateWeights` function, where multiple goroutines may attempt to update the `patternWeights` map concurrently without proper synchronization. This leads to inconsistent weights and inaccurate analysis.

To debug this, run the tool under high load and monitor the `patternWeights` map for inconsistencies. Using `go run -race ./cmd/log-analyzer` will help find this issue.

## Limitations

* The time window increase is not implemented.
* The `updateWeights` function is a stub.
* More robust burst handling is needed.
* Error messages are printed to standard error, not logged to a file.
* No debug mode is implemented.

## Notes

* The tool is designed to handle log streams from standard input.
* The tool is optimized for performance and efficiency.
* The tool is designed to be robust and handle various input scenarios.