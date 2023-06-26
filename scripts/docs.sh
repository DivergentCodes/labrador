#!/bin/bash

# https://pkg.go.dev/cmd/doc

for go_pkg in \
    "github.com/divergentcodes/labrador/cmd/labrador" \
    "github.com/divergentcodes/labrador/internal/aws" \
    "github.com/divergentcodes/labrador/internal/core" \
    "github.com/divergentcodes/labrador/internal/record" \
    ;
do
    printf "\n\n###############################################################################"
    printf "\n# Generating docs for: $go_pkg"
    printf "\n###############################################################################\n\n"
    go doc -all "$go_pkg"
done
