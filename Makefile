ROOTPROJECT ?= ../root
include ${ROOTPROJECT}/include.mk

.PHONY: clean build test
clean: go/clean
build: go/build
test: go/test
