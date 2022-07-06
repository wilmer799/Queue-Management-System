#!/bin/bash
echo 
go build fwq_visitor.go
./fwq_visitor 192.168.68.111 9093 8081 192.168.68.110 9092

