grep received_bytes ~/tmp/daemon.out  | jq .metric.value | uniq -c
