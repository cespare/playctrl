CWD=$(shell pwd)
export GOPATH=$(CWD)
SRC=$(shell find . -type f -name '*.go')
JS_SRC=$(shell find extension -type f -name '*.json' -o -name '*.coffee')

.PHONY: all clean fmt js

all: bin/playctrld bin/playctrl js

clean:
	rm -rf bin/*
	rm -rf ext_compiled/*

bin/playctrld: $(SRC)
	go build -o bin/playctrld server/server.go

bin/playctrl: $(SRC)
	go build -o bin/playctrl client/client.go

js: $(JS_SRC)
	rm -rf ext_compiled && mkdir ext_compiled
	cp extension/*.json extension/*.coffee ext_compiled/
	coffee -c ext_compiled/*.coffee
	rm -rf ext_compiled/*.coffee

fmt: $(SRC)
	@gofmt -s -l -w $(SRC)
