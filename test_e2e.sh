go test ./e2e
if [[ $? -ne 0 ]]; then
    echo "e2e tests failed" >&2
    exit 1
fi