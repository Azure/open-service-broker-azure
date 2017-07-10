#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

export REGISTRY=microsoft/

docker login -u "${DOCKER_USERNAME}" -p "${DOCKER_PASSWORD}"

if [[ "${TRAVIS_TAG}" =~ ^v[0-9]+\.[0-9]+\.[0-9]+[a-z]*$ ]]; then
    echo "Pushing images with tags '${TRAVIS_TAG}' and 'latest'."
    VERSION="${TRAVIS_TAG}" MUTABLE_TAG="latest" make docker-push
elif [[ "${TRAVIS_BRANCH}" == "master" ]]; then
    echo "Pushing images with default tags (git sha and 'canary')."
    make docker-push
else
    echo "Skipping deployment from non-master branch"
fi
