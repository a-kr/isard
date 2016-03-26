BIN=monitor

.PHONY: $(BIN) run clean fmt

export GOPATH:=$(shell pwd)
export CWD:=$(shell pwd)

$(BIN):
	go install $(BIN)

run: $(BIN)
	bin/$(BIN)

clean:
	rm -f bin/$(BIN)

fmt:
	go fmt $(BIN)

