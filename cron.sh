#!/bin/sh
rm -fr ./testssl.json
for n in $(cat domains.txt )
do
    docker run --rm -v `pwd`:/mnt drwetter/testssl.sh --append --jsonfile=/mnt/out.json $n
done
