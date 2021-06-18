GOOS=$(shell go env GOOS)
COMMIT_TAG=$(shell git rev-parse --short=10 HEAD)
COMMIT_TIME=$(shell git log --pretty=format:%cd --date=iso HEAD -1)
BUILD_TIME=$(shell date '+%Y-%m-%d %H:%M:%S')
GO_VERSION=$(shell go version)
VERSION=1.1

LDFLAGS= -X main.Version=$(VERSION)
LDFLAGS+= -X 'main.GoVersion=$(GO_VERSION)'
LDFLAGS+= -X main.CommitTag=$(COMMIT_TAG)
LDFLAGS+= -X 'main.BuildTime=$(BUILD_TIME)'
LDFLAGS+= -X 'main.CommitTime=$(COMMIT_TIME)'

default: build
generate:
	go generate -v ./...
asn.go:
	goyacc -o asnparser/asn.go -p "ParsedGrammar" asnparser/asn.y
build:
	go build -ldflags "$(LDFLAGS)" -o asngo cmd/asn1go/main.go
test:
	go test --count 1 -v ./uper/...
clean:
	rm -rf asnparser/asn.go
	rm -rf asngo
	rm -rf y.output
install:
	cp asngo /usr/local/bin/
uninstall:
	rm -rf /usr/local/bin/asngo
deps:
	go get golang.org/x/tools/cmd/goyacc