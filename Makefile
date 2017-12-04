.DEFAULT_GOAL := build

LD_FLAGS = -ldflags "-X main.VERSION=$(VERSION)"

all: build

build: build-deps
	CGO_ENABLED=0 GOOS=linux govendor build $(LD_FLAGS) -o bin/olowek

test: deps
	govendor test +l

build-deps: deps test
	@mkdir -p bin/

deps:
	@which govendor > /dev/null || \
	(go get -u github.com/kardianos/govendor)

clean:
	@rm -rf bin
.PHONY: all bump build release
