#!/usr/bin/env bash

curl \
  -v \
  -H "Accept: application/json" \
  -H "Content-Type: application/json" \
  -X POST \
  --data '{ "url": "https://github.com/cimple-ci/cimple-ruby-example.git", "commit": "master" }' \
  http://localhost:8080/builds
