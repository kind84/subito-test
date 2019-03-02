#!/bin/sh
# Build script
cd ../go/src/github.com/kind84/subito-test
CGO_ENABLED=0 GOOS=linux go build -a -v -installsuffix cgo