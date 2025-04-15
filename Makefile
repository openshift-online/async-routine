APP_NAME=asyncroutinemanager

MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
PROJECT_PATH := $(patsubst %/,%,$(dir $(MKFILE_PATH)))
LOCAL_BIN_PATH := $(PROJECT_PATH)/bin

export PATH := $(LOCAL_BIN_PATH):$(PATH)

.PHONY: all build test clean

all: build

MOCKGEN := $(LOCAL_BIN_PATH)/mockgen

.PHONY: mockgen-install
mockgen-install:
	@GOBIN=$(LOCAL_BIN_PATH) go install go.uber.org/mock/mockgen@v0.2.0 ;\

ginkgo:
	go install github.com/onsi/ginkgo/v2/ginkgo

tools: ginkgo


.PHONY: generate
generate: mockgen-install
	go generate ./...

build:
	go build -o $(APP_NAME)

test: tools generate
	ginkgo -r -v

clean:
	rm -f $(APP_NAME)
