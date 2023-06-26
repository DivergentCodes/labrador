#!/bin/bash

# Use goreleaser to build a snapshot of the binary for the current arch/os.
# Clean the ./dist folder before building.
# https://goreleaser.com/cmd/goreleaser_build/

goreleaser build --snapshot --clean --single-target

