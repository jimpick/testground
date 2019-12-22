#! /bin/bash

NUM_INSTANCES=$1

if [ -z "$NUM_INSTANCES" ]; then
  NUM_INSTANCES=2
fi

./testground -vv \
  run \
  --runner=cluster:k8s \
  --build-cfg bypass_cache=true \
  --build-cfg push_registry=true \
  --build-cfg registry_type=dockerhub \
  --instances $NUM_INSTANCES \
  basic-tcp/upload
