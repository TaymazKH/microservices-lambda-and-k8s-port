#!/bin/bash -eu

GOOS=linux
GOARCH=amd64
go build -o bootstrap .
zip -r deployment.zip bootstrap base.html
