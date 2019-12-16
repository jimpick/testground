#! /bin/bash

set +x

docker ps | sed -n 's,^.*\(tg-.*\)$,\1,p' | xargs docker stop
docker container prune -f
docker network prune -f
docker stop testground-sidecar || true
docker rm testground-sidecar || true
