#!/bin/bash

set -e

# Wait for the apiserver to start responding
./wupiao.sh broker 8080 300

exec bash -c mocha
