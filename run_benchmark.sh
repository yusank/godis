#!/bin/bash

# step 0: remove old files
rm -rf benchmark*.csv

# step 1: check params
clientNum=$1
benchNum=$2
if [ ! -n "$1" ] ;then
	clientNum=100
fi

if [ ! -n "$2" ] ;then
	benchNum=10000
fi

echo "clientNum: $clientNum"
echo "benchNum: $benchNum"

# step 2: run
day=$(date +%Y%m%d)
filePath=benchmark/benchmark.$day.csv
redis-benchmark -p 7379 -c $clientNum -n $benchNum -q --csv > $filePath

# step 3: open file
open $filePath
