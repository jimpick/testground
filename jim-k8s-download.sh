#! /bin/bash

for p in `kubectl get pods | awk '{ print $1 }' | grep tg-`; do
  kubectl logs $p
done
