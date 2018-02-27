# Checks for the existence of a docker client and prints a nice error message
# if it isn't present
.PHONY: check-docker
check-docker:
	@if [ -z $$(which docker) ]; then \
		echo "Missing \`docker\` client which is required for development"; \
		exit 2; \
	fi

# Checks for the existence of docker-compose and prints a nice error message if
# it isn't present
.PHONY: check-docker-compose
check-docker-compose: check-docker
	@if [ -z $$(which docker-compose) ]; then \
		echo "Missing \`docker-compose\` which is required for development"; \
		exit 2; \
	fi

.PHONY: pull-tile-generator
pull-tile-generator: check-docker
	docker pull cfplatformeng/tile-generator:v11.0.5

SHELL = /bin/bash

.PHONY: zip-broker
zip-broker: 
	pushd ../../..; \
	  rm -f contrib/cf/pcf-tile/resources/open-service-broker-azure.zip; \
	  zip -r contrib/cf/pcf-tile/resources/open-service-broker-azure.zip cmd pkg vendor; \
	popd;

.PHONY: prepare-to-generate-tile
prepare-to-generate-tile: pull-tile-generator check-docker-compose zip-broker

.PHONY: generate-major-tile
generate-major-tile: prepare-to-generate-tile
	docker-compose run --rm generate-major-tile

.PHONY: generate-minor-tile
generate-minor-tile: prepare-to-generate-tile
	docker-compose run --rm generate-minor-tile

.PHONY: generate-tile
generate-tile: prepare-to-generate-tile
	docker-compose run --rm generate-tile
