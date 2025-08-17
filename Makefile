# Copyright 2025 Sam Caldwell
# SPDX-License-Identifier: MIT

.PHONY: build lint test release

build:
	go build ./...

lint:
	go fmt ./...
	go vet ./...

test:
	go test ./...

release:
	go build -ldflags "-s -w" -o bin/mdlint cmd/mdlint/main.go
