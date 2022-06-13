#!/bin/bash
echo 
go build fwq_waitingTimeServer.go
./fwq_waitingTimeServer localhost 9094 localhost 9092

