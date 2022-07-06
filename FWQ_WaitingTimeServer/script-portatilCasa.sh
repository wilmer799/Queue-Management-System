#!/bin/bash
echo 
go build fwq_waitingTimeServer.go
./fwq_waitingTimeServer 192.168.68.110 9094 192.168.68.110 9092

