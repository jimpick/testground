#! /bin/bash

./testground -vv \
  run \
  --build-cfg bypass_cache=true \
  shaping-experiment/shape-traffic
