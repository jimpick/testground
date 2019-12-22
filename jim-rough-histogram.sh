#! /bin/bash

#grep time_to_receive ~/tmp/daemon.out | cut -c63- | jq .value | awk '{h[$1]++}END{for(i in h){print h[i],i|"sort -rn|head -40"}}' |awk '!max{max=$1;}{r="";i=s=60*$1/max;while(i-->0)r=r"#";printf "%15s %5d %s %s",$2,$1,r,"\n";}' | sort -rnk1
grep time_to_receive ~/tmp/daemon.out | cut -c1- | jq .metric.value | awk '{h[$1]++}END{for(i in h){print h[i],i|"sort -rn|head -40"}}' |awk '!max{max=$1;}{r="";i=s=60*$1/max;while(i-->0)r=r"#";printf "%15s %5d %s %s",$2,$1,r,"\n";}' | sort -rnk1
