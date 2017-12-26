#! /bin/bash

rm -rf build

# Build for 64bit.
GOARCH=amd64

GOOS=linux
go build -o build/linux/cLC

GOOS=windows
go build -o build/windows/cLC.exe

GOOS=darwin
go build -o build/darwin/cLC
