#!/bin/sh
touch ./cron.lock
rm -fr ./testssl.json
for n in $(cat domains.txt )
do
    docker run --rm -v `pwd`:/mnt drwetter/testssl.sh --append --jsonfile=/mnt/testssl.json $n
done
rm -fr ./cron.lock