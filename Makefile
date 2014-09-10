deps:
	go get github.com/tools/godep
	go get bitbucket.org/liamstask/goose/cmd/goose
	godep restore

install: deps

test: export GO_ENV=test
test:
	godep go test -v