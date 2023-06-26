#!/bin/bash

# Install pkgsite and serve docs for the project packages at http://localhost:8080.
# https://pkg.go.dev/golang.org/x/pkgsite/cmd/pkgsite

go install golang.org/x/pkgsite/cmd/pkgsite@latest

pkgsite
