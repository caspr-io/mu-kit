ROOTPROJECT ?= ../root
APIPROJECT = .
include ${ROOTPROJECT}/include.mk
PROTOC_FILES=river/sample.pb.go

.PHONY: clean build test generate
generate: ${PROTOC_FILES}
clean: go/clean
build: generate go/build
test: generate go/test

# Targets for cluster/up and cluster/teardown
.PHONY: up down
up:
down:
