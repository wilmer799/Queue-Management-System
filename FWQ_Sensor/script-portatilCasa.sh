#!/bin/bash
echo
go build fwq_sensor.go
./fwq_sensor 192.168.68.110 9092 atraccion15

