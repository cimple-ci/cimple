.PHONY: build
REVISION=`git rev-parse HEAD`
VERSION='0.0.1'
BUILD=`date +%FT%T%z`

build:
	go build -ldflags "-X main.VERSION=${VERSION} -X main.BuildDate=${BUILD} -X main.Revision=${REVISION}" -o output/cimple main.go
