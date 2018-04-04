SHELL = /bin/bash

GENERATE_TILE = docker run --rm -v ${CURDIR}:/workspace -w /workspace --name=tile-generator cfplatformeng/tile-generator:v11.0.5 tile build

# Checks for the existence of a docker client and prints a nice error message
# if it isn't present
.PHONY: check-docker
check-docker:
	@if [ -z $$(which docker) ]; then \
		echo "Missing \`docker\` client which is required for development"; \
		exit 2; \
	fi

.PHONY: pull-tile-generator
pull-tile-generator: check-docker
	docker pull cfplatformeng/tile-generator:v11.0.5

.PHONY: zip-broker
zip-broker:
	pushd ../../..; \
	  rm -f contrib/cf/pcf-tile/resources/open-service-broker-azure.zip; \
	  zip -r contrib/cf/pcf-tile/resources/open-service-broker-azure.zip cmd pkg vendor; \
	popd;

.PHONY: prepare-to-generate-tile
prepare-to-generate-tile: pull-tile-generator zip-broker

.PHONY: generate-major-tile
generate-major-tile: prepare-to-generate-tile
	${GENERATE_TILE} major

.PHONY: generate-minor-tile
generate-minor-tile: prepare-to-generate-tile
	${GENERATE_TILE} minor

.PHONY: generate-tile
generate-tile: prepare-to-generate-tile
	${GENERATE_TILE}
