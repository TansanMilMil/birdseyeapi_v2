#!/bin/bash -eu

cd `dirname $0`
BASE_DIR="go"
GO_FILE="$BASE_DIR/src/main.go"

go run $GO_FILE
