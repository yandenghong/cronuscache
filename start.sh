#!/bin/bash
trap "rm server;kill 0" EXIT

go build -o server
./server -port=8309 &
./server -port=8310 &
./server -port=8311 -api=1 &

sleep 2
echo ">>> start test"
curl "http://localhost:9999/api?key=Tom" &
curl "http://localhost:9999/api?key=Tom" &
curl "http://localhost:9999/api?key=Tom" &

wait
