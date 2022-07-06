#!/bin/bash
echo 
go build fwq_engine.go
./fwq_engine 192.168.68.110 9092 3 5,6,4,9 192.168.68.110 9094

