#!/bin/bash

export ES_SERVER=127.0.0.1:9200

LISTEN_ADDRESS=10.29.1.1:12345 STORAGE_ROOT=/tmp/1/ go run ../internal/dataServer/dataServer.go &
LISTEN_ADDRESS=10.29.1.2:12345 STORAGE_ROOT=/tmp/2/ go run ../internal/dataServer/dataServer.go &
LISTEN_ADDRESS=10.29.1.3:12345 STORAGE_ROOT=/tmp/3/ go run ../internal/dataServer/dataServer.go &
LISTEN_ADDRESS=10.29.1.4:12345 STORAGE_ROOT=/tmp/4/ go run ../internal/dataServer/dataServer.go &
LISTEN_ADDRESS=10.29.1.5:12345 STORAGE_ROOT=/tmp/5/ go run ../internal/dataServer/dataServer.go &
LISTEN_ADDRESS=10.29.1.6:12345 STORAGE_ROOT=/tmp/6/ go run ../internal/dataServer/dataServer.go &

go run ../internal/apiServer
