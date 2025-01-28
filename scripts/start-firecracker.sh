#!/bin/bash

SOCKET_PATH=$1
LOG_FILE=$2

if [ -z "$SOCKET_PATH" ] || [ -z "$LOG_FILE" ]; then
  echo "Usage: $0 <socket_path> <log_file>"
  exit 1
fi

if [ -e "$SOCKET_PATH" ]; then
  rm -f "$SOCKET_PATH"
fi

firecracker --api-sock $SOCKET_PATH > $LOG_FILE 2>&1 &
