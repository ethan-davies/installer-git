#!/bin/bash

mkdir -p bin

GOOS=linux GOARCH=amd64 go build -o bin/setup-linux
GOOS=darwin GOARCH=amd64 go build -o bin/setup-macos
GOOS=windows GOARCH=amd64 go build -o bin/setup-windows.exe

echo "Finished building setup binarys"