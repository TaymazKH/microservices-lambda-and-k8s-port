#!/bin/bash -eu

GOOS=linux
GOARCH=amd64
go build -o service .
zip -r deployment.zip bootstrap handler.sh service
