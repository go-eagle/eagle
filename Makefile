SHELL := /bin/bash
BASEDIR = $(shell pwd)

# build with version infos
versionDir = "github.com/go-eagle/eagle/pkg/version"
gitTag = $(shell if [ "`git describe --tags --abbrev=0 2>/dev/null`" != "" ];then git describe --tags --abbrev=0; else git log --pretty=format:'%h' -n 1; fi)
buildDate = $(shell TZ=Asia/Shanghai date +%FT%T%z)
gitCommit = $(shell git log --pretty=format:'%H' -n 1)
gitTreeState = $(shell if git status|grep -q 'clean';then echo clean; else echo dirty; fi)

ldflags="-w -X ${versionDir}.gitTag=${gitTag} -X ${versionDir}.buildDate=${buildDate} -X ${versionDir}.gitCommit=${gitCommit} -X ${versionDir}.gitTreeState=${gitTreeState}"

PROJECT_NAME := "github.com/go-eagle/eagle"
PKG := "$(PROJECT_NAME)"
GO_VERSION=$(shell go version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f2)
PKG_LIST := $(shell go list ${PKG}/... | grep -v examples | grep -v pkg)
GO_FILES := $(shell find . -name '*.go' | grep -v _test.go)

# init environment variables
export PATH        := $(shell go env GOPATH)/bin:$(PATH)
export GOPATH      := $(shell go env GOPATH)
export GO111MODULE := on

# make   make all
.PHONY: all
all: lint test build

.PHONY: build
# make build, Build the binary file
build: dep
	@go build -v -ldflags ${ldflags} .

.PHONY: dep
# make dep Get the dependencies
dep:
	@go mod download

.PHONY: fmt
# make fmt
fmt:
	@gofmt -s -w .

.PHONY: golint
# make golint
golint:
	@if ! which golint &>/dev/null; then \
  		echo "Installing golint"; \
  		go get -u golang.org/x/lint/golint; \
  	fi
	@golint -set_exit_status ${PKG_LIST}

.PHONY: lint
# make lint
lint:
	@if ! which golangci-lint &>/dev/null; then \
  		echo "Installing golangci-lint"; \
  		go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.43.0; \
  	fi
	${GOPATH}/bin/golangci-lint run ./...

.PHONY: test
# make test
test:
	@go test -short ${PKG_LIST}

.PHONY: vet
# make vet
vet:
	go vet ./... | grep -v vendor;true

.PHONY: cover
# make cover
cover:
	@go test -short -coverprofile coverage.txt -covermode=atomic ${PKG_LIST}

.PHONY: view-cover
# make view-cover  preview coverage
view-cover:
	go tool cover -html=coverage.txt

.PHONY: docker
# make docker  生成docker镜像
docker:
	docker build -t eagle:v1.0 -f Dockerfile .

.PHONY: clean
# make clean
clean:
	@-rm -vrf eagle
	@-rm -vrf cover.out
	@-rm -vrf coverage.txt
	@go mod tidy
	@echo "clean finished"

.PHONY: docs
# gen swagger doc
docs:
	@if ! which swag &>/dev/null; then \
  		echo "downloading swag"; \
  		go get -u github.com/swaggo/swag/cmd/swag; \
  	fi
	@swag init
	@echo "gen-docs done"
	@echo "see docs by: http://localhost:8080/swagger/index.html"

.PHONY: graph
# make graph 生成交互式的可视化Go程序调用图(会在浏览器自动打开)
graph:
	@export GO111MODULE="on"
	@if ! which go-callvis &>/dev/null; then \
  		echo "downloading go-callvis"; \
  		go get -u github.com/ofabry/go-callvis; \
  	fi
	@echo "generating graph"
	@go-callvis github.com/go-eagle/eagle

.PHONY: mockgen
# make mockgen gen mock file
mockgen:
	@echo "downloading mockgen"
	@go get github.com/golang/mock/mockgen
	cd ./internal &&  for file in `egrep -rnl "type.*?interface" ./dao | grep -v "_test" `; do \
		echo $$file ; \
		cd .. && mockgen -destination="./internal/mock/$$file" -source="./internal/$$file" && cd ./internal ; \
	done

.PHONY: ca
# generate ca
ca:
	openssl req -new -nodes -x509 -out config/server.crt -keyout config/server.key -days 3650 -subj "/C=DE/ST=NRW/L=Earth/O=Random Company/OU=IT/CN=127.0.0.1/emailAddress=xxxxx@qq.com"

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := all

