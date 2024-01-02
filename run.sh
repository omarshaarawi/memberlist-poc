#!/bin/bash

# Configuration
START_PORT=8000
START_LISTEN_PORT=8080
NUM_INSTANCES=10
EXECUTABLE="./memberlist-poc"
SLEEP_TIME=2

# Start the first node separately
FIRST_NODE_PORT=$START_PORT
FIRST_NODE_NAME="node-$FIRST_NODE_PORT"

# Start the first instance
# echo the command to be executed
echo $EXECUTABLE -port $FIRST_NODE_PORT -listen ":$START_LISTEN_PORT" -name $FIRST_NODE_NAME
$EXECUTABLE -port $FIRST_NODE_PORT -listen ":$START_LISTEN_PORT" -name $FIRST_NODE_NAME &
echo "Started $FIRST_NODE_NAME"
sleep $SLEEP_TIME

# Subsequent nodes will join the first node
PEERS="127.0.0.1:$FIRST_NODE_PORT"

# Loop to start remaining instances
for (( i=1; i<$NUM_INSTANCES; i++ ))
do
    PORT=$((START_PORT + i))
    LISTEN_PORT=$((START_LISTEN_PORT + i))
    NAME="node-$PORT"

    # Start the instance in the background
    # echo the command to be executed
    echo $EXECUTABLE -port $PORT -listen ":$LISTEN_PORT" -name $NAME -peers $PEERS
    $EXECUTABLE -port $PORT -listen ":$LISTEN_PORT" -name $NAME -peers $PEERS &
    
    echo "Started $NAME"
    sleep $SLEEP_TIME

    # Update PEERS to include the new instance for subsequent instances
    PEERS="$PEERS,127.0.0.1:$PORT"
done

echo "Launched $NUM_INSTANCES instances."

