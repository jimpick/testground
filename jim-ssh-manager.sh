#! /bin/bash

eval `grep manager_dns infra/docker-swarm/ansible/inventories/jim3`

echo $manager_dns

ssh ubuntu@$manager_dns
