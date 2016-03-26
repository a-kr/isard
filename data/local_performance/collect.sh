#!/bin/bash -e

NOW=$(date "+%Y-%m-%d %H:%M:%S")

LA=$(cat /proc/loadavg | cut -d ' ' -f 1)
USEDMEM=$(free -m | grep 'buffers/cache' | awk '{print $3}')
NPROC=$(ps aux | wc -l)

echo $NOW,$LA,$USEDMEM,$NPROC
