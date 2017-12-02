#!/usr/bin/env bash
# [w]ait [u]ntil [p]ort [i]s [a]ctually [o]pen
# Usage: ./wupiao.sh server:host tries

if [ -z "$1" ] || [ -z "$2" ] || [ -z "$3" ]; then echo "Usage: wupiao <ip> <port> <timeout in seconds>"; exit 1; fi

set -e

function echo_red {
  echo -e "\033[0;31m$1\033[0m"
}

IP=$1
PORT=$2
TIMEOUT=$3

COUNTER=1

until nc -z $IP $PORT &> /dev/null; do
  if [ $COUNTER -gt $TIMEOUT ]; then
    echo_red "Timed out waiting for $IP:$PORT"
    exit 1
  fi
  sleep 1
  let COUNTER=COUNTER+1
done;