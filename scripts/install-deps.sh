#!/bin/bash

DOCKER_COMPOSE_VERSION=1.16.1
DOCKER_COMPOSE_LOCATION=$(which docker-compose)
DOCKER_COMPOSE_DL_URL="https://github.com/docker/compose/releases/download/${DOCKER_COMPOSE_VERSION}/docker-compose-Linux-x86_64"
sudo curl -o ${DOCKER_COMPOSE_LOCATION} -L ${DOCKER_COMPOSE_DL_URL}
chmod +x ${DOCKER_COMPOSE_LOCATION}
docker version
docker-compose version
