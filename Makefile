PROG=bin/proxy


SRCS=./cmd

# git commit hash
COMMIT_HASH=$(shell git rev-parse --short HEAD || echo "GitNotFound")

# date
BUILD_DATE=$(shell date '+%Y-%m-%d %H:%M:%S')

# flag
CFLAGS = -ldflags "-s -w -X \"main.BuildVersion=${COMMIT_HASH}\" -X \"main.BuildDate=$(BUILD_DATE)\""

all:
	if [ ! -d "./bin/" ]; then \
	mkdir bin; \
	fi
	go build $(CFLAGS) -o $(PROG) $(SRCS)

# race version
race:
	if [ ! -d "./bin/" ]; then \
    	mkdir bin; \
    	fi
	go build $(CFLAGS) -race -o $(PROG) $(SRCS)

# remote debug
DLVFLAGS = -ldflags "-X \"main.BuildVersion=${COMMIT_HASH}\" -X \"main.BuildDate=$(BUILD_DATE)\""
DLVGCFLAGS = -gcflags "all=-N -l"
dlv:
	if [ ! -d "./bin/" ]; then \
	mkdir bin; \
	fi
	go build $(DLVFLAGS) $(DLVGCFLAGS) -o $(PROG) $(SRCS)

clean:
	rm -rf ./bin

run:
	go run .

run_dev:
	go run .

test:
	go test -v --v ./...