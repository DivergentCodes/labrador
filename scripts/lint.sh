#!/bin/bash

# Lint the Go code.

# https://pkg.go.dev/cmd/vet
# Vet examines Go source code and reports suspicious constructs.
# It can find errors not caught by the compilers.
go vet ./...

# https://staticcheck.io/docs/
# Using static analysis, it finds bugs and performance issues, offers
# simplifications, and enforces style rules.
go install honnef.co/go/tools/cmd/staticcheck@latest
staticcheck ./...
