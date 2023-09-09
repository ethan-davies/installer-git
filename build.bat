@echo off

mkdir bin

set GOOS=linux
set GOARCH=amd64
go build -o bin\setup-linux

set GOOS=darwin
set GOARCH=amd64
go build -o bin\setup-macos

set GOOS=windows
set GOARCH=amd64
go build -o bin\setup-windows.exe

echo Finished building setup binaries
