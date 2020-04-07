#!/bin/bash

docker run -d --name influxdb-1 -p 8086:8086 influxdb:1.7
docker run -d --name influxdb-2 -p 8087:8086 influxdb:1.7
