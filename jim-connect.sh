#! /bin/bash

eval `grep manager_dns infra/docker-swarm/ansible/inventories/jim3`

echo $manager_dns

ssh -nNT -L 4545:/var/run/docker.sock ubuntu@$manager_dns
