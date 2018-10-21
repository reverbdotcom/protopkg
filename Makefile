.PHONY: all test

APPNAME := protopkg

all: install test
install:
	@go install

test:
	@go test

release: bin/$(APPNAME).linux bin/$(APPNAME).darwin

bin/$(APPNAME).linux:
	@env GOOS=linux go build -o $@

bin/$(APPNAME).darwin:
	@env GOOS=darwin go build -o $@
