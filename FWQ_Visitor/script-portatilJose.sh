#!/bin/bash
echo 
go build fwq_visitor.go
./fwq_visitor 192.168.43.50 9093 8081 192.168.43.247 9092

