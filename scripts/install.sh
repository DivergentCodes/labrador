#!/bin/bash

# Locally "install" the application by copying the binary to the GOPATH.
# Get the binary name from the GoReleaser configuration file.

binary_name=$(grep "project_name" .goreleaser.yaml | sed 's/project_name:[[:space:]]\+//')
echo "BIN:   [$binary_name]"

os=$(go env GOOS)
arch=$(go env GOARCH)
echo "OS:    [$os]"
echo "ARCH:  [$arch]"

gopath=$(go env GOPATH)
src_path="$(find ./dist -name $binary_name)"
dst_path="$gopath/bin/$binary_name"
echo "SRC:   [$src_path]"
echo "DST:   [$dst_path]"

cp "$src_path" "$dst_path"
