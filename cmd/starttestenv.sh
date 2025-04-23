#!/bin/bash

export ES_SERVER=127.0.0.1:9200

prefix=~/storage_data

LISTEN_ADDRESS=10.29.1.1:12345 STORAGE_ROOT=${prefix}/1/ go run ../internal/dataServer/dataServer.go > ../log/dataserver1.log &
LISTEN_ADDRESS=10.29.1.2:12345 STORAGE_ROOT=${prefix}/2/ go run ../internal/dataServer/dataServer.go > ../log/dataserver2.log &
LISTEN_ADDRESS=10.29.1.3:12345 STORAGE_ROOT=${prefix}/3/ go run ../internal/dataServer/dataServer.go > ../log/dataserver3.log &
LISTEN_ADDRESS=10.29.1.4:12345 STORAGE_ROOT=${prefix}/4/ go run ../internal/dataServer/dataServer.go > ../log/dataserver4.log &
LISTEN_ADDRESS=10.29.1.5:12345 STORAGE_ROOT=${prefix}/5/ go run ../internal/dataServer/dataServer.go > ../log/dataserver5.log &
LISTEN_ADDRESS=10.29.1.6:12345 STORAGE_ROOT=${prefix}/6/ go run ../internal/dataServer/dataServer.go > ../log/dataserver6.log &

#go run ../myes > ../log/es.log 2>&1 &

go run ../internal/apiServer > ../log/apiserver.log 2>&1 &
