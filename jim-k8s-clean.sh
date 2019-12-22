#! /bin/bash

set +x

kubectl get jobs | awk '{ print $1 }' | grep tg- | xargs kubectl delete job
kubectl get pods | awk '{ print $1 }' | grep tg- | xargs kubectl delete pod 
