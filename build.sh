#!/bin/bash
binName="thoughts_export"

GOOS=windows GOARCH=amd64 go build -o "bin/$binName"_win
GOOS=linux GOARCH=amd64 go build -o "bin/$binName"_linux
GOOS=darwin GOARCH=amd64 go build -o "bin/$binName"_macos

echo "build working is done, see below binary files"

ls "bin/$binName"*
