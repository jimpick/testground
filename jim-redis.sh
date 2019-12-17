#! /bin/bash

echo REDIS_HOST=localhost:6379
redis-server --notify-keyspace-events "\$szxK" --save "" --appendonly no
