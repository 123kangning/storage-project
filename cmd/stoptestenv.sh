#!/bin/bash

# 停止 apiServer
apiServer_pids=$(pgrep -f "apiServer")
if [[ -n "$apiServer_pids" ]]; then
    echo "Stopping apiServer ..."
    for pid in $apiServer_pids; do
        kill "$pid"
    done
fi

# 停止 dataServer
dataServer_pids=$(pgrep -f "dataServer")
if [[ -n "$dataServer_pids" ]]; then
    echo "Stopping dataServer ..."
    for pid in $dataServer_pids; do
        kill "$pid"
    done
fi
