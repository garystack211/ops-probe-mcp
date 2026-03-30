#!/bin/bash
mkdir -p ./bin
GOOS=darwin GOARCH=arm64 go build -o ./bin/ops-probe-mcp-darwin-arm64
GOOS=linux  GOARCH=amd64 go build -o ./bin/ops-probe-mcp-linux-amd64
