SHELL := /bin/bash
BASEDIR = $(shell pwd)

# build with verison infos
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

.PHONY: all
all: build

.PHONY: build
build: ## Build the binary file
	@go build -v -ldflags ${ldflags} .

.PHONY: docker
# make docker
# 生成docker镜像
docker:
	docker build \
		-t snake:$(versionDir) \
		-f Dockeffile \
		.

.PHONY: clean
clean:
	rm -f snake
	rm cover.out coverage.txt
	find . -name "[._]*.s[a-w][a-z]" | xargs -i rm -f {}

.PHONY: gotool
gotool:
	gofmt -w .
	go tool vet . | grep -v vendor;true

.PHONY: dep
dep: ## Get the dependencies
	@go mod download

.PHONY: lint
lint: ## Lint Golang files
	@golint -set_exit_status ${PKG_LIST}

.PHONY: ci-lint-prepare
ci-lint-prepare:
	@echo "Installing golangci-lint"
    @curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s latest

.PHONY: ci-lint
ci-lint:
	@./bin/golangci-lint run ./...

.PHONY: test
test: ## Run unittests
	@go test -short ${PKG_LIST}

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@go test -short -coverprofile cover.out -covermode=atomic ${PKG_LIST}
	@cat cover.out >> coverage.txt

.PHONY: test-view
# view test result
test-view:
	@go tool cover -html=coverage.txt

.PHONY: gen-docs
gen-docs:
	@swag init
	@mv docs/docs.go api/http
	@mv docs/swagger.json api/http
	@mv docs/swagger.yaml api/http
	@echo "gen-docs done"
	@echo "see docs by: http://localhost:8080/swagger/index.html"

.PHONY: gen-graph
# 生成交互式的可视化Go程序调用图
gen-graph:
	@echo "downloading go-callvis"
	@go get github.com/1024casts/snake
	@echo "generating graph"
	@go-callvis github.com/1024casts/snake

.PHONY: ca
ca:
	openssl req -new -nodes -x509 -out conf/server.crt -keyout conf/server.key -days 3650 -subj "/C=DE/ST=NRW/L=Earth/O=Random Company/OU=IT/CN=127.0.0.1/emailAddress=xxxxx@qq.com"



