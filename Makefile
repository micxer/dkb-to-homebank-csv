TEST?=./...

.PHONY: default help tidy build clean tests

BIN_NAME=dkb2homebank

default: build

help:
	@echo 'Management commands for dkb2homebank:'
	@echo
	@echo 'Usage:'
	@echo '    make build           Compile the project.'
	@echo '    make tidy            Runs go mod tidy, mostly used for ci.'
	@echo '    make test            Run tests on a compiled project.'
	@echo '    make clean           Clean the directory tree.'
	@echo

build:
	@echo "building ${BIN_NAME}"
	@echo "GOPATH=${GOPATH}"
	go test
	go build -o ${BIN_NAME} converter.go

tidy:
	go mod tidy

clean:
	@test ! -e bin/${BIN_NAME} || rm bin/${BIN_NAME}

# test: runs the unit tests
test:
	go test
