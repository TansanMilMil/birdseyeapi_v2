#!/bin/bash -eu

cd $(dirname $0)

go test ./go/src/...
