-include .env

compile:
	@go build -o calendar ./cmd/calendar/main.go
	@go build -o sender ./cmd/sender/main.go
	@go build -o scheduler ./cmd/scheduler/main.go

install:
	@go install ./...