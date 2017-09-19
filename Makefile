VERSION ?= $(shell git describe --always --abbrev=7 --dirty)

BINARY_DIR := bin
BINARY_NAME := asb

BASE_IMAGE_NAME      = azure-service-broker
MUTABLE_TAG         ?= canary
IMAGE_NAME           = $(REGISTRY)$(BASE_IMAGE_NAME):$(VERSION)
MUTABLE_IMAGE_NAME   = $(REGISTRY)$(BASE_IMAGE_NAME):$(MUTABLE_TAG)

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

# Deletes any existing asb binary AND destroys any running containers AND
# destroys the dev environment image.
.PHONY: clean
clean: check-docker-compose
	rm -rf ${BINARY_DIR}
	rm -rf ${CONTRIB_BINARY_DIR}
	docker-compose down --rmi local &> /dev/null

# Containerized project bootstrapping-- requires docker-compose
# This will (re)build the development environment and populate the vendor/
# directory with dependencies specified by glide.lock
.PHONY: dev-bootstrap
dev-bootstrap: check-docker-compose
	docker-compose build dev
	docker-compose run --rm dev glide install

# Containerized dependency update-- requires docker-compose
# This will (re)build the development environment, populate the vendor/
# directory with updated dependencies, and update glide.lock accordingly 
.PHONY: dev-update
dev-update: check-docker-compose
	docker-compose build dev
	docker-compose run --rm dev glide up

# Allow developers to step into the containerized development environment--
# requires docker-compose
.PHONY: dev
dev: check-docker-compose
	docker-compose run --rm dev bash

.PHONY: test
test: test-unit test-module-lifecycles

# Containerized unit tests-- requires docker-compose
.PHONY: test-unit
test-unit: check-docker-compose
	docker-compose run --rm test bash -c 'go test -tags unit $$(glide nv)'

# Containerized module lifecycle tests-- requires docker-compose
.PHONY: test-module-lifecycles
test-module-lifecycles: check-docker-compose
	@docker-compose run \
		--rm \
		-e AZURE_SUBSCRIPTION_ID=${AZURE_SUBSCRIPTION_ID} \
		-e AZURE_TENANT_ID=${AZURE_TENANT_ID} \
		-e AZURE_CLIENT_ID=${AZURE_CLIENT_ID} \
		-e AZURE_CLIENT_SECRET=${AZURE_CLIENT_SECRET} \
		test \
		bash -c 'go test \
			-parallel 10 \
		  -timeout 60m \
			github.com/Azure/azure-service-broker/tests/lifecycle -v'

.PHONY: lint
lint: check-docker-compose
	docker-compose run \
		--rm test \
		bash -c 'gometalinter $$(glide nv) \
			--disable-all \
			--enable gofmt \
			--enable vet \
			--enable vetshadow \
			--enable gotype \
			--enable deadcode \
			--enable golint \
			--enable varcheck \
			--enable structcheck \
			--enable aligncheck \
			--enable errcheck \
			--enable megacheck \
			--enable ineffassign \
			--enable interfacer \
			--enable unconvert \
			--enable goconst \
			--enable gas \
			--enable goimports \
			--enable misspell \
			--enable unparam \
			--enable lll \
			--line-length 80 \
			--deadline 120s'

# Running the tests starts a containerized Redis dedicated to testing (if it
# isn't already running). It's left running afterwards (to speed up the next
# execution). It remains running unless explicitly shut down. This is a
# convenience task for stopping it.
.PHONY: stop-test-redis
stop-test-redis: check-docker-compose
	docker-compose kill test-redis
	docker-compose rm -f test-redis

# Containerized binary build for linux/64 only-- requires docker-compose
.PHONY: build
build: check-docker-compose
	docker-compose run --rm dev \
		go build -o ${BINARY_DIR}/${BINARY_NAME} ./cmd/broker

# (Re)Build the Docker image for the asb and run it
.PHONY: run
run: check-docker-compose build
	@# Force the docker-compose "broker" service to be rebuilt-- this is separate
	@# from the docker-build task used to produce a correctly tagged Docker image,
	@# although both builds are based on the same Dockerfile
	docker-compose build broker
	docker-compose run \
		--rm -p 8080:8080 \
		-e AZURE_SUBSCRIPTION_ID=${AZURE_SUBSCRIPTION_ID} \
		-e AZURE_TENANT_ID=${AZURE_TENANT_ID} \
		-e AZURE_CLIENT_ID=${AZURE_CLIENT_ID} \
		-e AZURE_CLIENT_SECRET=${AZURE_CLIENT_SECRET} \
		broker

# Running the broker starts a containerized Redis dedicated to that purpose (if
# it isn't already running). It's left running afterwards (to speed up the next
# execution AND to retain state so broker recovery from incomplete async
# operations can be demonstrated). It remains running unless explicitly shut
# down. This is a convenience task for stopping it.
.PHONY: stop-broker-redis
stop-broker-redis: check-docker-compose
	docker-compose kill broker-redis
	docker-compose rm -f broker-redis

# Build the Docker image
.PHONY: docker-build
docker-build: check-docker build
	docker build -t $(IMAGE_NAME) .
	docker tag $(IMAGE_NAME) $(MUTABLE_IMAGE_NAME)

# Push the Docker image
.PHONY: docker-push
docker-push: check-docker docker-build
	docker push $(IMAGE_NAME)
	docker push $(MUTABLE_IMAGE_NAME)

# ---------------------------------------------------------------------------- #
# contrib/                                                                     #
# ---------------------------------------------------------------------------- #

CONTRIB_BINARY_DIR := contrib/bin
CLI_BINARY_NAME := broker-cli

.PHONY: build-mac-broker-cli
build-mac-broker-cli: check-docker-compose
	docker-compose run --rm -e GOOS=darwin -e GOARCH=amd64 dev \
		go build -o ${CONTRIB_BINARY_DIR}/${CLI_BINARY_NAME} \
		./contrib/cmd/cli
