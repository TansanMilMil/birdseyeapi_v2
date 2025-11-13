#!/bin/bash -eu

cd `dirname $0`
BASE_DIR="go"
DIST_DIR="$BASE_DIR/dist"
BIN_PATH="$DIST_DIR/birdseyeapi_v2"
GO_FILE="$BASE_DIR/src/main.go"


mkdir -p $DIST_DIR
go build -o $BIN_PATH $GO_FILE

echo "binary built! -> $BIN_PATH"
