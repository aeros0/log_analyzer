#!/bin/bash

# Configuration variables
LOG_FILE="test_logs.log"
MIN_LOGS_PER_SECOND=600  # Minimum log timestamps per log-time second
MAX_LOGS_PER_SECOND=4900 # Maximum log timestamps per log-time second
MIN_SLEEP=0.01          # Minimum system-time sleep (in seconds)
MAX_SLEEP=0.3           # Maximum system-time sleep (in seconds)

# Array of possible error messages
declare -a ERROR_MESSAGES=(
    "Database connection failed"
    "Null pointer exception"
    "File not found"
    "Access denied"
    "Out of memory"
    "Network timeout occurred"
    "Illegal argument provided"
    "User authentication failed"
)

# Function to generate a random IP address in the format 192.168.X.Y
generate_ip() {
    echo "192.168.$((RANDOM % 254 + 1)).$((RANDOM % 254 + 1))"
}

# Function to get current timestamp in ISO 8601 format
get_timestamp() {
    date -u +"%Y-%m-%dT%H:%M:%SZ"
}

# Function to generate a random log level
generate_level() {
    local rand=$((RANDOM % 3))
    case $rand in
        0) echo "ERROR" ;;
        1) echo "INFO" ;;
        2) echo "DEBUG" ;;
    esac
}

# Function to generate an error message for ERROR level logs
generate_error_message() {
    local level=$1
    if [ "$level" = "ERROR" ]; then
        local index=$((RANDOM % ${#ERROR_MESSAGES[@]}))
        echo "Error 500 - ${ERROR_MESSAGES[$index]}"
    else
        echo ""
    fi
}

# Function to generate a single log entry with a given timestamp
generate_log_entry() {
    local level=$(generate_level)
    local ip=$(generate_ip)
    local error_message=$(generate_error_message "$level")
    echo "[$1] $level - IP:$ip $error_message" # Use provided timestamp
}

# Create or clear the log file
> "$LOG_FILE"

echo "Starting log generation. Press Ctrl+C to stop..."

# Main loop to generate logs
while true; do
    # Generate a random number of logs per log-time second
    logs_per_second=$((RANDOM % (MAX_LOGS_PER_SECOND - MIN_LOGS_PER_SECOND + 1) + MIN_LOGS_PER_SECOND))

    # Get the initial timestamp for the "log-time" second
    initial_timestamp=$(get_timestamp)

    # Generate log entries for the "log-time" second
    for ((i=0; i<logs_per_second; i++)); do
        generate_log_entry "$initial_timestamp" >> "$LOG_FILE"
    done

    # Random sleep between "log-time" seconds (in system-time seconds)
    sleep_time=$(echo "scale=2; $MIN_SLEEP + (rand() / ${RANDOM_MAX:-32767}) * ($MAX_SLEEP - $MIN_SLEEP)" | bc -l)
    sleep "$sleep_time"
done