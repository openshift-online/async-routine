APP_NAME=asyncroutinemanager

.PHONY: all build test clean

all: build

ginkgo:
	go install github.com/onsi/ginkgo/v2/ginkgo

tools: ginkgo

build:
	go build -o $(APP_NAME)

test: tools
	ginkgo -r -v

clean:
	rm -f $(APP_NAME)
