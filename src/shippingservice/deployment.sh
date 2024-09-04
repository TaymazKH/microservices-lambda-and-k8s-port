#!/bin/bash -eu

GOOS=linux
GOARCH=amd64
go build -o service ./main.go
zip -r deployment.zip bootstrap handler.sh service  # todo: deploy with SAM
