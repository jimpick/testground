#! /bin/bash

# sudo apt install gnupg2 pass
# $(aws ecr get-login --no-include-email --region us-west-2)

docker tag ipfs/testground:latest 909427826938.dkr.ecr.us-west-2.amazonaws.com/testground:latest

docker push 909427826938.dkr.ecr.us-west-2.amazonaws.com/testground:latest

# Doesn't work - Docker global services aren't privileged
# docker service create --name testground-sidecar --env REDIS_HOST=testground-redis --network control --mode global --constraint "node.labels.TGRole == worker" --with-registry-auth 909427826938.dkr.ecr.us-west-2.amazonaws.com/testground:latest sidecar --runner docker

# Shell
# docker run -it --name testground-sidecar --cap-add NET_ADMIN --cap-add SYS_ADMIN --pid host --env REDIS_HOST=testground-redis --network control --mount type=bind,source=/var/run/docker.sock,destination=/var/run/docker.sock --entrypoint sh 909427826938.dkr.ecr.us-west-2.amazonaws.com/testground:latest

# Run sidecar
echo ----------------------------------------
echo sudo snap install aws-cli --classic
echo sudo apt install gnupg2 pass
echo '$(aws ecr get-login --no-include-email --region us-west-2)'
echo 'docker stop testground-sidecar; docker rm testground-sidecar'
echo docker run --name testground-sidecar --detach --cap-add NET_ADMIN --cap-add SYS_ADMIN --pid host --env REDIS_HOST=testground-redis --network control --mount type=bind,source=/var/run/docker.sock,destination=/var/run/docker.sock 909427826938.dkr.ecr.us-west-2.amazonaws.com/testground:latest sidecar --runner docker

