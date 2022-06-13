#!/bin/bash
echo
go build fwq_sensor.go
./fwq_sensor localhost 9092 atraccion15

