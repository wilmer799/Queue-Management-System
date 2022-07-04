#!/bin/bash
echo 
go build fwq_engine.go
./fwq_engine 192.168.43.247 9092 2 192.168.43.247 9094

