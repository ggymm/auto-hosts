#!/bin/sh

cd ../ || exit

echo "build windows"
export GOOS=windows
export GOARCH=amd64
go build -o auto-hosts.exe