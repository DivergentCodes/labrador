#!/bin/bash

# Use goreleaser to create a snapshot release across all configured arch/os combinations.
# Clean the ./dist folder before building.
# https://goreleaser.com/cmd/goreleaser_release/

# A GoReleaser run is split into 4 major steps:
#     defaulting: configures sensible defaults for each step
#     building: builds the binaries, archives, packages, Docker images, etc
#     publishing: publishes the release to the configured SCM, Docker registries, blob storages...
#     announcing: announces your release to the configured channels

goreleaser release --clean
