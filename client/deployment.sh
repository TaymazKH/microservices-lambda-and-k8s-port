#!/bin/bash -eu

GOOS=linux
GOARCH=amd64
go build -o service ./client.go
zip -r deployment.zip bootstrap handler.sh service
