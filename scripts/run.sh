#!/bin/sh
# Run script
cd ../go/src/github.com/kind84/subito-test
go install
cd $GOPATH/bin
subito-test