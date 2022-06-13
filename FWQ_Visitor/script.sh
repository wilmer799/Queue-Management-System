#!/bin/bash
echo 
go build fwq_visitor.go
./fwq_visitor localhost 9093 8081 localhost 9092

