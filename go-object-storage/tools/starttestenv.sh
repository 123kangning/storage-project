#!/bin/bash

export RABBITMQ_SERVER="amqp://kangning:9264wkn.@localhost:5672/vhost-1"
export ES_SERVER=localhost:9200

LISTEN_ADDRESS=10.29.1.1:12345 STORAGE_ROOT=/tmp/1 go run ../final/dataServer/dataServer.go &
LISTEN_ADDRESS=10.29.1.2:12345 STORAGE_ROOT=/tmp/2 go run ../final/dataServer/dataServer.go &
LISTEN_ADDRESS=10.29.1.3:12345 STORAGE_ROOT=/tmp/3 go run ../final/dataServer/dataServer.go &
LISTEN_ADDRESS=10.29.1.4:12345 STORAGE_ROOT=/tmp/4 go run ../final/dataServer/dataServer.go &
LISTEN_ADDRESS=10.29.1.5:12345 STORAGE_ROOT=/tmp/5 go run ../final/dataServer/dataServer.go &
LISTEN_ADDRESS=10.29.1.6:12345 STORAGE_ROOT=/tmp/6 go run ../final/dataServer/dataServer.go &

LISTEN_ADDRESS=10.29.2.1:12345 go run ../final/apiServer/apiServer.go &
LISTEN_ADDRESS=10.29.2.2:12345 go run ../final/apiServer/apiServer.go &

WEB_SERVER_ADDRESS=10.29.3.1:12345 API_SERVER=10.29.2.1:12345 go run ../special/webServer/webServer.go
