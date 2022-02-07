#!/bin/sh
CGO_ENABLED=0 go build -o termora -v -ldflags="-buildid= -X github.com/termora/berry/common.Version=`git rev-parse --short HEAD`" ./cmd
