#!/bin/bash
statik -f -src=./web/dist/

GOOS=linux

GOARCH=amd64

go build -o goploy main.go

GOOS=darwin

GOARCH=amd64

go build -o goploy.mac main.go

GOOS=windows

GOARCH=amd64

go build -o goploy.exe main.go
