GOARCH = amd64
OPTS = CGO_ENABLED=0

build:
	go build -v -o ko cmd/ko/main.go

windows.goos: EXT = .exe

%.goos:
	$(OPTS) GOOS=$* go build -ldflags="-s -w" -v -o ko-$*-$(GOARCH)$(EXT) cmd/ko/main.go

.PHONY: all
all: windows.goos linux.goos darwin.goos

run:
	go run cmd/ko/main.go
