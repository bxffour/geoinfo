#!/bin/bash

set -e

CGO_ENABLED=0 go build -o ./bin/bs
docker build -t ghcr.io/bxffour/crest/bootstrap:1.1.0 .