#!/bin/bash

# Reset the environment
./e2e/reset.sh

# Clear the test cache (as Go doesn't always pick up the fact that the source files have changed, when running the E2E tests)
go clean --testcache

# Run the E2E Go tests, and handle test failure
go test ./e2e
if [[ $? -ne 0 ]]; then
    echo "e2e tests failed" >&2
    exit 1
fi