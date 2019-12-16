#! /bin/bash

./testground -vv \
  run \
  --build-cfg bypass_cache=true \
  basic-tcp/upload
