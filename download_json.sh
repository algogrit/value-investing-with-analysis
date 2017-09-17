#!/usr/bin/env bash

./to-screener-url.sh $1 | xargs curl | jq '.' > data/$1.json
