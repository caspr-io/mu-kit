ROOTPROJECT ?= ../root
APIPROJECT = .
PROTOBUF_FILES=streaming/sample.pb.go
include ${ROOTPROJECT}/include.mk

.PHONY: clean build test
clean: go/clean
build: protobuf/generate go/build
test: protobuf/generate go/test

# Targets for cluster/up and cluster/teardown
.PHONY: up down
up:
down:
