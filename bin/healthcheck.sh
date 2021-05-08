#!/bin/bash

set -euxo "pipefail"
IFS=$'\n\t'

curl \
    --fail \
    "curl localhost:5003/time/"
