#!/bin/bash
echo 
go build fwq_engine.go
./fwq_engine localhost 9092 3 localhost 9094

