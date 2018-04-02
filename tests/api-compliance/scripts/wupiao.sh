#!/usr/bin/env bash
# [w]ait [u]ntil [p]ort [i]s [a]ctually [o]pen
# Usage: ./wupiao.sh host port timeout

if [ -z "$1" ] || [ -z "$2" ] || [ -z "$3" ]; then echo "Usage: ./wupiao <host> <port> <timeout in seconds>"; exit 1; fi

set -eo pipefail

function echo_red {
  echo -e "\033[0;31m$1\033[0m"
}

HOST=$1
PORT=$2
TIMEOUT=$3

COUNTER=1

until nc -z $HOST $PORT &> /dev/null; do
  if [ $COUNTER -gt $TIMEOUT ]; then
    echo_red "Timed out waiting for $HOST:$PORT"
    exit 1
  fi
  sleep 1
  let COUNTER=COUNTER+1
done;