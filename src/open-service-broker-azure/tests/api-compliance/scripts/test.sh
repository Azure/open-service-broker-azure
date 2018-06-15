#!/bin/bash
# Usage: ./test-api-compliance.sh host port timeout

################################################################################
# This is a COPY of the test script included in docker-osb-checker. It exists  #
# here because certain constraints of CircleCI do not permit us to use that    #
# image to conduct compliance testing. This is described in further detail     #
# below where local modifications to this script are grouped and explained.    #
################################################################################

if [ -z "$1" ] || [ -z "$2" ] || [ -z "$3" ]; then echo "Usage: ./test-api-compliance.sh <host> <port> <timeout in seconds>"; exit 1; fi

set -eo pipefail

HOST=$1
PORT=$2
TIMEOUT=$3

TESTS_PATH=/opt/osb-checker/2.13/tests
CONFIG_PATH=/$TESTS_PATH/test/configs/config_mock.json

################################################################################
# Begin workaround...                                                          #
################################################################################

# To run compliance tests in CircleCI, we need to compile and run go code (a
# dummy broker) AND then execute some mocha tests (nodejs). Circle will want the
# tests to be in the "primary" container. This is the container in which all
# subsequent test steps are executed. The node container would be "primary" in
# that case and the go container would be "secondary." Since no steps can be
# specified to run in a secondary container-- including the checkout step--
# there would no way to get the go code we want to compile and run into that
# secondary container.
#
# Our workaround is to do everything in a single container. On principle, I
# (krancour), refuse to bundle node tools into our containerized GO development
# environment. The resulting images would simply be too big and would be useful
# only to a small subset of projects that utilize both go and node. So, for this
# one application, in this one place, we will install node and the OSB checker
# just-in-time into the running, PRIMARY go container. We can live with this
# short term whilst we avoid consuming as-of-late unreliable Circle machine
# executors, whilst knowing that a more elegant solution to this exists if/when
# we move to Jenkins.
#
# If the workaround seems to do some inefficient things, its deliberate for the
# sake of keeping these edits isolated to one spot in the script and easily
# reversed if and when the constraints that led us here are ever lifted.

curl -sL https://deb.nodesource.com/setup_9.x | bash -
apt-get update
apt-get install -y nodejs netcat
npm install -g mocha
git clone https://github.com/openservicebrokerapi/osb-checker.git /opt/osb-checker
rm $TESTS_PATH/test/configs/*.json
cp /go/src/github.com/Azure/open-service-broker-azure/tests/api-compliance/localhost-config.json /app/config.json
cd $TESTS_PATH
npm install

################################################################################
# End workaround. There were no more modifications to the script past this     #
# point.                                                                       #
################################################################################

ln -s /app/config.json $CONFIG_PATH || true

# Wait for the apiserver to start responding
/app/wupiao.sh $HOST $PORT $TIMEOUT

cd $TESTS_PATH
exec bash -c mocha
