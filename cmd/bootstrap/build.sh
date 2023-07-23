#!/bin/bash

set -e

CGO_ENABLED=0 go build -o ./bin/bs
docker build -t ghcr.io/bxffour/geoinfo/bootstrap:0.9.2 .