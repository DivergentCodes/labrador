#!/bin/bash

# Use goreleaser to create a snapshot release across all configured arch/os combinations.
# Clean the ./dist folder before building.
# https://goreleaser.com/cmd/goreleaser_release/
goreleaser release --snapshot --clean
