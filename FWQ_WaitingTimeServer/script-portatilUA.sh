#!/bin/bash
echo 
go build fwq_waitingTimeServer.go
./fwq_waitingTimeServer 192.168.43.247 9094 192.168.43.247 9092

