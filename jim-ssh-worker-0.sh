#! /bin/bash

WORKER=`grep 'worker\.0' -A1 infra/docker-swarm/ansible/inventories/jim3 | tail -1`

set +x

ssh ubuntu@$WORKER
