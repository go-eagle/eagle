SHELL := /bin/bash
BASEDIR = $(shell pwd)

# build with version infos
versionDir = "github.com/1024casts/snake/pkg/version"
gitTag = $(shell if [ "`git describe --tags --abbrev=0 2>/dev/null`" != "" ];then git describe --tags --abbrev=0; else git log --pretty=format:'%h' -n 1; fi)
buildDate = $(shell TZ=Asia/Shanghai date +%FT%T%z)
gitCommit = $(shell git log --pretty=format:'%H' -n 1)
gitTreeState = $(shell if git status|grep -q 'clean';then echo clean; else echo dirty; fi)

ldflags="-w -X ${versionDir}.gitTag=${gitTag} -X ${versionDir}.buildDate=${buildDate} -X ${versionDir}.gitCommit=${gitCommit} -X ${versionDir}.gitTreeState=${gitTreeState}"

PROJECT_NAME := "github.com/1024casts/snake"
PKG := "$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

.PHONY: build
# make build, Build the binary file
build:
	go build -v -ldflags ${ldflags} .

.PHONY: dep
# make dep
dep: ## Get the dependencies
	@go mod download

.PHONY: fmt
# make fmt
fmt:
	@gofmt -s -w .

.PHONY: lint
# make lint
lint:
	@golint -set_exit_status ${PKG_LIST}

.PHONY: ci-lint
# make ci-lint
ci-lint: ci-lint-prepare
	@./bin/golangci-lint run ./...

.PHONY: ci-lint-prepare
# make ci-lint-prepare
ci-lint-prepare:
	@echo "Installing golangci-lint"
    @curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s latest

.PHONY: test
# make test
test: test-case vet-case
	@go test -short ${PKG_LIST}

.PHONY: test-case
# make test-case
test-case:
	@go test -cover ./... | grep -v vendor;true

.PHONY: vet-case
# make vet-case
vet-case:
	go vet ./... | grep -v vendor;true

.PHONY: gen-coverage
# make gen-coverage
gen-coverage:
	@go test -short -coverprofile cover.out -covermode=atomic ${PKG_LIST}
	@cat cover.out >> coverage.txt

.PHONY: review-cover
# make review-cover
review-cover:
	@go tool cover -html=coverage.txt


.PHONY: docker
# make docker  生成docker镜像
docker:
	docker build -t snake:$(versionDir) -f Dockeffile .

.PHONY: clean
# make clean
clean:
	@-rm -vrf snake
	@-rm -vrf cover.out
	@-rm -vrf coverage.txt
	@go mod tidy
	@echo "clean finished"

.PHONY: docs
# gen swagger doc
docs:
	@swag init
	@mv docs/docs.go api/http
	@mv docs/swagger.json api/http
	@mv docs/swagger.yaml api/http
	@echo "gen-docs done"
	@echo "see docs by: http://localhost:8080/swagger/index.html"

.PHONY: graph
# 生成交互式的可视化Go程序调用图
graph:
	@export GO111MODULE="on"
	@if ! which go-callvis &>/dev/null; then \
  		echo "downloading go-callvis"; \
  		go get github.com/ofabry/go-callvis; \
  	fi
	@go get github.com/ofabry/go-callvis
	@echo "generating graph"
	@go-callvis github.com/1024casts/snake

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

.DEFAULT_GOAL := help

