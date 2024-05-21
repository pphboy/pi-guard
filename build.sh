# 编译到linux平台
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o pidns
