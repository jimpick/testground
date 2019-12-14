#! /bin/bash

./testground -vv \
  run \
  --build-cfg bypass_cache=true \
  smlbench2/simple-add-get
