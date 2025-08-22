APP := awsvault
BIN := bin/$(APP)

.PHONY: all
all: fmt build

.PHONY: build
build:
	go build -v -o $(BIN)

.PHONY: install-deps
install-deps:
	go get .

.PHONY: clean
clean:
	rm -rf bin

.PHONY: fmt
fmt:
	go fmt ./...
