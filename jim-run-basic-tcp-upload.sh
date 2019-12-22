#! /bin/bash

NUM_INSTANCES=$1

if [ -z "$NUM_INSTANCES" ]; then
  NUM_INSTANCES=2
fi

./testground -vv \
  run \
  --build-cfg bypass_cache=true \
  --instances $NUM_INSTANCES \
  basic-tcp/upload
