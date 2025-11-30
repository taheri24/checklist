BINARY := checklist
CMD := ./cmd/checklist
OUTPUT := bin/$(BINARY)

.PHONY: all build install test clean

all: build

build:
@mkdir -p bin
GO111MODULE=on go build -o $(OUTPUT) $(CMD)

install:
GO111MODULE=on go install $(CMD)

test:
GO111MODULE=on go test ./...

clean:
rm -rf bin
