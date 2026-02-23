#!/bin/bash

# Check for fswatch
if ! command -v fswatch >/dev/null 2>&1; then
    echo "Error: fswatch not found. Install it with: brew install fswatch"
    exit 1
fi

# Check for .go files
go_files=$(find . -name '*.go')
if [ -z "$go_files" ]; then
    echo "Error: No .go files found in current directory"
    exit 1
fi

# Ensure _build directory exists
mkdir -p _build

# Log the files being monitored
echo "Monitoring the following .go files:"
echo "$go_files"

while true; do
    # Check for stop condition
    if [ -f "stop.txt" ]; then
        echo "Stop file detected. Exiting."
        pkill -f '_build/aclips' || true
        rm stop.txt
        exit 0
    fi

    # Build the Go program
    echo "Building at $(date)..."
    go build -o _build/aclips
    if [ $? -ne 0 ]; then
        echo "Build failed at $(date)" >> build_errors.log
        sleep 1
        continue
    fi

    # Kill any existing aclips process
    echo "Killing existing _build/aclips process..."
    pkill -f '_build/aclips' || echo "No _build/aclips process was running"

    # Start the new aclips process in the background
    echo "Starting _build/aclips..."
    ./_build/aclips &
    aclips_pid=$!
    echo $aclips_pid > aclips.pid
    echo "Started _build/aclips with PID $aclips_pid"

    # Wait for changes to .go files using fswatch
    echo "Waiting for changes to .go files..."
    fswatch -0 -o $go_files | while read -r -d $'\0' event; do
        echo "File change detected at $(date): $event"
        break
    done
done
