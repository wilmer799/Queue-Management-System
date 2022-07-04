#!/bin/bash
echo 
go build fwq_waitingTimeServer.go
./fwq_waitingTimeServer 192.168.56.102 9094 192.168.56.102 9092

