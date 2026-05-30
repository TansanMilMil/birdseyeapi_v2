#!/bin/sh
set -e

sudo curl -sL https://taskfile.dev/install.sh | sudo sh -s -- -b /usr/local/bin
go mod download
curl -fsSL https://claude.ai/install.sh | bash
