#!/bin/bash

#mac 版的binary
go build -o build/refactor/mac/gdeyamlOperator cmd/refactor/main.go

#Linux版的binary
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/gdeyamlOperator cmd/refactor/main.go
