#!/bin/bash

apiServer_pid=$(pgrep -f "apiServer")
if [[ -n "$apiServer_pid" ]]; then
    echo "Stopping apiServer ..."
    kill $apiServer_pid
fi

dataServer_pid=$(pgrep -f "dataServer")
if [[ -n "$dataServer_pid" ]]; then
    echo "Stopping dataServer ..."
    kill $dataServer_pid
fi

webServer_pid=$(pgrep -f "webServer")
if [[ -n "$webServer_pid" ]]; then
    echo "Stopping webServer ..."
    kill $webServer_pid
fi
