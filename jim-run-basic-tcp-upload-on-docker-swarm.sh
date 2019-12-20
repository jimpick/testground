#! /bin/bash

./testground -vv \
  run \
  --runner=cluster:swarm \
  --build-cfg bypass_cache=true \
  --build-cfg push_registry=true \
  --build-cfg registry_type=aws \
  --run-cfg keep_service=true \
  basic-tcp/upload
