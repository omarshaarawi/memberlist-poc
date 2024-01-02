#!/bin/bash

# Configuration
START_PORT=8000
NUM_INSTANCES=10
LOG_FILES=()

# Generate list of log files
for (( i=0; i<$NUM_INSTANCES; i++ )); do
    PORT=$((START_PORT + i))
    LOG_FILES+=("node-$PORT.log")
done

# Launch multitail with all log files
multitail "${LOG_FILES[@]}"

