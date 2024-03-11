.SILENT:
.DEFAULT_GOAL := fast-run

include ./.env

export ENABLE_DEBUG_LOGS
export POM_BOT_TOKEN
export UPDATES_CHECK_PERIOD_SECS

.PHONY: build
build:
	go build -o ./build/pomd ./cmd/pomd/main.go

.PHONY: run
run:
	./build/pomd

.PHONY: fast-run
fast-run:
	go run ./cmd/pomd/main.go
