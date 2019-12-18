#! /bin/bash

./testground -vv \
  run \
  --runner=cluster:swarm \
  --build-cfg bypass_cache=true \
  --build-cfg push_registry=true \
  --build-cfg registry_type=aws \
  basic-tcp/upload
