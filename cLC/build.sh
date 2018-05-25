#! /bin/bash

rm -rf build

# Build for current platform
go build -o build/current/cLC -ldflags "-s -w"

# Make sure Windows users get an .exe
cp build/current/cLC build/current/cLC.exe

# Build for 64bit Linux and Windows
GOARCH=amd64 GOOS=linux go build -o build/linux/cLC -ldflags "-s -w"

GOARCH=amd64 GOOS=windows go build -o build/windows/cLC.exe -ldflags "-s -w"
