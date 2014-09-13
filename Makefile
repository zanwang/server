env ?= dev
ts = $(shell /bin/date "+%Y-%m-%d-%H-%M-%S")
tarball = build/$(ts).tar.bz2

deps:
	go get github.com/tools/godep
	godep restore

install: deps

test: export GO_ENV=test
test:
	godep go test -v

migrate:
	go get bitbucket.org/liamstask/goose/cmd/goose
	goose -env $(env) up

build:
	GOARCH=amd64 GOOS=linux go build -o majimoe
	mkdir -p build
	tar -jc -f $(tarball) --exclude build --exclude .git --exclude Godeps ./
	rm -f ./majimoe

.PHONY: build