#!/bin/bash

# https://pkg.go.dev/cmd/doc

for go_pkg in \
    "divergent.codes/labrador/cmd/labrador" \
    "divergent.codes/labrador/internal/aws" \
    "divergent.codes/labrador/internal/core" \
    "divergent.codes/labrador/internal/record" \
    ;
do
    printf "\n\n###############################################################################"
    printf "\n# Generating docs for: $go_pkg"
    printf "\n###############################################################################\n\n"
    go doc -all "$go_pkg"
done
