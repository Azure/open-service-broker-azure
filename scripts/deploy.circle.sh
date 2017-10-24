#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

# This variable is quite possibly undefined. If it is, define a safe default
# value (empty string) to avoid problems with the nounset option.
CIRCLE_TAG=${CIRCLE_TAG:-""}

docker login -u "${DOCKER_HUB_USERNAME}" -p "${DOCKER_HUB_PASSWORD}"

if [[ "${CIRCLE_TAG}" =~ ^v[0-9]+\.[0-9]+\.[0-9]+[a-z]*$ ]]; then
    echo "Pushing images with tags '${CIRCLE_TAG}' and 'latest'."
    REGISTRY=microsoft/ VERSION="${CIRCLE_TAG}" MUTABLE_TAG="latest" \
      make docker-push
elif [[ "${CIRCLE_BRANCH}" == "master" ]]; then
    echo "Pushing images with default tags (git sha and 'canary')."
    REGISTRY=microsoft/ make docker-push
else
    echo "Skipping deployment from non-master branch"
fi
