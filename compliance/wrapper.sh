#!/bin/bash
#wrapper.sh

set -e

until [[ 200 -eq $(curl --write-out %{http_code} --silent --output /dev/null http://broker:8080/healthz) ]]; do
    echo "Waiting for server to be up"
    sleep 5
done

exec bash -c mocha

