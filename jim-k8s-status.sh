#! /bin/bash

watch -n 5 "kubectl get pod | grep ^tg- | awk '{ print \$3 }' | sort | uniq -c"
