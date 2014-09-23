CWD=$(shell pwd)
export GOPATH=$(CWD)
SRC=$(shell find . -type f -name '*.go')
EXT_SRC=$(shell find extension -type f)

.PHONY: all clean fmt

all: bin/playctrl extension.zip

clean:
	rm -rf bin/*
	rm -rf extension.zip

bin/playctrl: $(SRC)
	go build -o bin/playctrl

extension.zip: extension
	zip -r extension.zip extension/*

fmt: $(SRC)
	@gofmt -s -l -w $(SRC)
