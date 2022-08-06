COMMIT_SHA_SHORT ?= $(shell git rev-parse --short=12 HEAD)
PWD_DIRR:= ${CURDIR}
SHELL := /bin/bash


default: help;

# ======================================================================================

build: ## build the debian package
	@GOOS=linux GOARCH=amd64 go build -o out/x1-acpi-manager main.go
	@nfpm package -f nfpm.yaml -p deb

publish: build ## publish the release to github
	@zarf/publish.sh

clean: ## clean
	@rm *.deb || true
	@rm out/x1-acpi-manager || true

help: ## Show this help
	@egrep '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST)  | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m·%-20s\033[0m %s\n", $$1, $$2}'