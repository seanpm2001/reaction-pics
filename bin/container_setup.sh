#!/bin/bash

set -euo pipefail
IFS=$'\n\t'

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"
cd "$DIR/.." || exit 1

# Compile go and node
mkdir /root/gocode/bin
mkdir /root/gocode/pkg
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
export PATH="$PATH:$GOPATH/bin"
make bins
npm prune --production
