#! /bin/bash

./testground -vv \
  run \
  --builder docker:go \
  --runner local:docker \
  --build-cfg bypass_cache=true \
  --test-param timeout_secs=300 \
  --instances=3 \
  dht/find-peers
