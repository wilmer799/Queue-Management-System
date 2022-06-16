#!/bin/bash
echo 
go build fwq_engine.go
./fwq_engine localhost 9092 3 7,6,4,8 localhost 9094

