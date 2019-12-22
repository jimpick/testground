#! /bin/bash

./testground -vv \
  run \
  --runner=cluster:k8s \
  --build-cfg bypass_cache=true \
  --build-cfg push_registry=true \
  --build-cfg registry_type=dockerhub \
  basic-tcp/upload
