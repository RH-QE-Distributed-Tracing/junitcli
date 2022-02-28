GO = $(shell which go)

build: ## Builds the CLI
	$(GO) build -o bin/junitcli cmd/junitcli/junitcli.go

install: ## Installs the utility
	$(GO) install cmd/junitcli/junitcli.go

help: ## Show this help
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'
