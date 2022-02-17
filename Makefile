GO = $(shell which go)

build:
	$(GO) build -o bin/junitcli cmd/junitcli/junitcli.go

install:
	$(GO) install cmd/junitcli/junitcli.go
