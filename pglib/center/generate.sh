#!/bin/sh
#
# 生成脚本
protoc --go_out=. --go_opt=paths=source_relative \
	   --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	   ct.proto

