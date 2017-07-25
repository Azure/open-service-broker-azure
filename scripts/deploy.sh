#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

# docker login -u "${DOCKER_USERNAME}" -p "${DOCKER_PASSWORD}"
docker login -u "${QUAY_USERNAME}" -p "${QUAY_PASSWORD}" quay.io

if [[ "${TRAVIS_TAG}" =~ ^v[0-9]+\.[0-9]+\.[0-9]+[a-z]*$ ]]; then
    echo "Pushing images with tags '${TRAVIS_TAG}' and 'latest'."
    REGISTRY=quay.io/deis/ VERSION="${TRAVIS_TAG}" MUTABLE_TAG="latest" \
      make docker-push
elif [[ "${TRAVIS_BRANCH}" == "master" ]]; then
    echo "Pushing images with default tags (git sha and 'canary')."
    REGISTRY=quay.io/deisci/ make docker-push
else
    echo "Skipping deployment from non-master branch"
fi
