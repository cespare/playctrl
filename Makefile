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

extension.zip: $(EXT_SRC)
	mkdir ext_compiled
	cp -r extension/*.json extension/*.coffee extension/icon ext_compiled/
	coffee -c ext_compiled/*.coffee
	rm -rf ext_compiled/*.coffee
	rm -rf ext_compiled/icon/*.xcf
	zip -r extension.zip ext_compiled/*
	rm -rf ext_compiled

fmt: $(SRC)
	@gofmt -s -l -w $(SRC)
