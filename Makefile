SHELL := /bin/bash
BASEDIR = $(shell pwd)

# build with verison infos
versionDir = "snake/pkg/version"
gitTag = $(shell if [ "`git describe --tags --abbrev=0 2>/dev/null`" != "" ];then git describe --tags --abbrev=0; else git log --pretty=format:'%h' -n 1; fi)
buildDate = $(shell TZ=Asia/Shanghai date +%FT%T%z)
gitCommit = $(shell git log --pretty=format:'%H' -n 1)
gitTreeState = $(shell if git status|grep -q 'clean';then echo clean; else echo dirty; fi)

ldflags="-w -X ${versionDir}.gitTag=${gitTag} -X ${versionDir}.buildDate=${buildDate} -X ${versionDir}.gitCommit=${gitCommit} -X ${versionDir}.gitTreeState=${gitTreeState}"

PROJECT_NAME := "github.com/1024casts/snake"
PKG := "$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

all: build
build: dep ## Build the binary file
	@go build -v -ldflags ${ldflags} .
clean:
	rm -f snake
	rm cover.out coverage.txt
	find . -name "[._]*.s[a-w][a-z]" | xargs -i rm -f {}
gotool:
	gofmt -w .
	go tool vet . | grep -v vendor;true
dep: ## Get the dependencies
	@go mod download
lint: ## Lint Golang files
	@golint -set_exit_status ${PKG_LIST}
test: ## Run unittests
	@go test -short ${PKG_LIST}
test-coverage: ## Run tests with coverage
	@go test -short -coverprofile cover.out -covermode=atomic ${PKG_LIST}
	@cat cover.out >> coverage.txt
test-view: ## view test result
	@go tool cover -html=coverage.txt
swag-init:
	swag init
	@echo "swag init done"
	@echo "see docs by: http://localhost:8080/swagger/index.html"

ca:
	openssl req -new -nodes -x509 -out conf/server.crt -keyout conf/server.key -days 3650 -subj "/C=DE/ST=NRW/L=Earth/O=Random Company/OU=IT/CN=127.0.0.1/emailAddress=xxxxx@qq.com"

help:
	@echo "make - compile the source code"
	@echo "make clean - remove binary file and vim swp files"
	@echo "make gotool - run go tool 'fmt' and 'vet'"
	@echo "make ca - generate ca files"
	@echo "make swag-init - gen swag doc"

.PHONY: clean gotool ca help


