build:
	go build -ldflags="-s -w" -v -o ko.exe cmd/ko/main.go

run:
	go run cmd/ko/main.go

default: build
