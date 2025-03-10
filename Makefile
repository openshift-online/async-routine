APP_NAME=asyncroutinemanager

.PHONY: all build test clean

all: build

build:
	go build -o $(APP_NAME)

test:
	ginkgo -r -v

clean:
	rm -f $(APP_NAME)
