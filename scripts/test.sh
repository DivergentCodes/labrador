#!/bin/bash

# Run tests, generate code coverage HTML report at ".coverage.html,"
# and check if it is above the defined minimum coverage (default 0%).
# Set COVERAGE_MIN to enforce code coverage minimum.

set -e -u

export CGO_ENABLED=0

# Run tests and generate code coverage.
go test ./... -coverprofile ".coverage.out" -covermode atomic

# Generate HTML report.
go tool cover -html=".coverage.out" -o ".coverage.html"

# Set default coverage minimum if not defined.
: ${COVERAGE_MIN:=00.0}

# Get coverage percentage.
coverage_percent=`go tool cover -func=".coverage.out" | tail -n 1 | sed -Ee 's!^[^[:digit:]]+([[:digit:]]+(\.[[:digit:]]+)?)%$!\1!'`
result=`echo "$coverage_percent >= $COVERAGE_MIN" | bc`

# Check coverage percentage against minimum.
test "$result" -eq 1 && exit 0
echo "Insufficient code coverage: $coverage_percent% (<$COVERAGE_MIN%)" >&2
exit 1
