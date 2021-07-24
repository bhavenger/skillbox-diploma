#!/bin/sh

cd /tests
go mod download
GO111MODULE=on CGO_ENABLED=0 GOOS=linux go test -v ./...