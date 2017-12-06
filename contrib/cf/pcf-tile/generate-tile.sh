#!/bin/bash
pushd ../../..
  rm -f contrib/cf/pcf-tile/resources/open-service-broker-azure.zip
  zip -r contrib/cf/pcf-tile/resources/open-service-broker-azure.zip cmd pkg vendor
popd
if [ "$1" = "-major" ]; then
  tile build major
elif [ "$1" = "-minor" ]; then
  tile build minor
else
  tile build
fi
