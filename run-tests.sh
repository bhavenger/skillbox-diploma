#!/bin/bash

docker run \
    -v $(pwd):/tests \
    golang:1.16.6-alpine3.14 \
    /tests/scripts/test_in_docker.sh
