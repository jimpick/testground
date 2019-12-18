#! /bin/bash

set +x

./testground -vv \
  run \
  --build-cfg bypass_cache=true \
  basic-tcp/upload-with-shaping
