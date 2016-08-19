.PHONY: all deps test build clean

REPO=github.com/svenwltr/i3-statusbar
NAME=$(notdir $(REPO))

VERSION=$(shell git describe --always --dirty | tr '-' '.' )
BUILDOPTS=-ldflags "-X main.version=$(VERSION)" \

usage:
	@echo "USAGE: make [all] [bootstrap] [clean] [test] [build] [install] [cov] [cc]"
	
all: bootstrap test build cov cc

bootstrap: tools deps

clean:
	rm -rf target vendor

deps:
ifneq ($(strip $(shell go list)),$(REPO))
	$(error This repository should be cloned into $$GOPATH/src/$(REPO))
endif
	glide install

test: deps
	go test -cover $(shell glide nv)

build: deps
	mkdir -p target
	go build \
		$(BUILDOPTS) \
		-o target/$(NAME)

install: test
	go install $(BUILDOPTS)

tools:
	go get gopkg.in/matm/v1/gocov-html
	go get github.com/axw/gocov/gocov

cov: deps tools
	gocov test $(shell glide nv) \
		| gocov-html > target/coverage.html

cc: test
	mkdir -p target
	GOOS=linux GOARCH=amd64 go build $(BUILD_OPTS) -o target/$(NAME)-linux-amd64
	GOOS=darwin GOARCH=amd64 go build $(BUILD_OPTS) -o target/$(NAME)-darwin-amd64
