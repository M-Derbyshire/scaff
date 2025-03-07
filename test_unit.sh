#!/bin/bash

go test $(go list ./... | grep -v /e2e) # Run all tests, excluding the e2e directory
if [[ $? -ne 0 ]]; then
    echo "unit tests failed" >&2
    exit 1
fi