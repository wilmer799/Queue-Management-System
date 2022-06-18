#!/bin/bash
echo 
go build fwq_engine.go
./fwq_engine localhost 9092 3 5,6,4,9 localhost 9094

