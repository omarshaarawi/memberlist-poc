# memberlist-poc

## Overview
This repository contains a proof of concept for a distributed system using `memberlist` library in Go.

The PoC demonstrates the creation of a simple HTTP server that interacts with a `memberlist` cluster. Nodes in the cluster can join, leave, and broadcast messages to other nodes. It also includes scripts for starting multiple instances and monitoring logs.

## Structure
- `main.go`: The main Go file implementing the HTTP server and memberlist node.
- `run.sh`: Bash script to start multiple instances of the memberlist node.
- `log.sh`: Bash script for aggregating and viewing logs from multiple instances.

## Requirements
- Go (version 1.15 or later recommended)
- `memberlist` library by Hashicorp
- `multitail` (optional, for log monitoring)

## Usage
- To start nodes: Run `./run.sh`. It starts a specified number of node instances.
- To view logs: Run `./log.sh`. It aggregates logs from all running instances for easier monitoring.

## Configuration
- `START_PORT`: Starting port number for the first node (Default: 8000).
- `START_LISTEN_PORT`: Starting listen port (Default: 8080).
- `NUM_INSTANCES`: Number of node instances to launch (Default: 10).
- `SLEEP_TIME`: Time to wait between launching instances (Default: 2 seconds).

