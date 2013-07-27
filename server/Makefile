EXE=bin/playctrld
CWD=$(shell pwd)
export GOPATH=$(CWD)
PROJECT_SRC=$(shell find . -path ./vendor -prune -o -type f -name '*.go' -print)
ALL_SRC=$(shell find . -type f -name '*.go')

.PHONY: all clean run deps fmt

all: $(EXE)

clean:
	rm -rf $(EXE)

run: $(EXE)
	./$(EXE)

deps:
	git submodule update --init

$(EXE): $(ALL_SRC) deps
	go build -o $(EXE)

fmt: $(PROJECT_SRC)
	@gofmt -s -l -w $(PROJECT_SRC)
