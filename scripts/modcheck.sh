#!/bin/bash

# https://go.dev/ref/mod#go-mod-verify
# Compare package hashes to verify they haven't been changed since download.
go mod verify
