GIT_VERSION = $(shell git describe --always --abbrev=7 --dirty)

BINARY_DIR := bin
BINARY_NAME := osba

BASE_IMAGE_NAME        = azure-open-service-broker

RC_IMAGE_NAME          = $(REGISTRY)$(BASE_IMAGE_NAME):$(GIT_VERSION)
RC_MUTABLE_IMAGE_NAME  = $(REGISTRY)$(BASE_IMAGE_NAME):canary

REL_IMAGE_NAME         = $(REGISTRY)$(BASE_IMAGE_NAME):$(REL_VERSION)
REL_MUTABLE_IMAGE_NAME = $(REGISTRY)$(BASE_IMAGE_NAME):latest

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

# Checks to ensure that AZURE_* environment variables needed to run the broker
# and its integration tests are set
.PHONY: check-azure-env-vars
check-azure-env-vars:
ifndef AZURE_SUBSCRIPTION_ID
	$(error AZURE_SUBSCRIPTION_ID is not defined)
endif
ifndef AZURE_TENANT_ID
	$(error AZURE_TENANT_ID is not defined)
endif
ifndef AZURE_CLIENT_ID
	$(error AZURE_CLIENT_ID is not defined)
endif
ifndef AZURE_CLIENT_SECRET
	$(error AZURE_CLIENT_SECRET is not defined)
endif

# Deletes any existing osba binary AND destroys any running containers AND
# destroys the dev environment image.
.PHONY: clean
clean: check-docker-compose
	rm -rf ${BINARY_DIR}
	rm -rf ${CONTRIB_BINARY_DIR}
	docker-compose down --rmi local &> /dev/null

# Allow developers to step into the containerized development environment--
# requires docker-compose
.PHONY: dev
dev: check-docker-compose
	docker-compose run --rm dev bash

# Containerized dependency install/update-- requires docker-compose
.PHONY: dep
dep: check-docker-compose
	docker-compose run --rm dev dep ensure -v

.PHONY: verify-vendored-code
verify-vendored-code: check-docker-compose
	docker-compose run --rm dev bash -c ' \
		export PRJ_DIR=$$(pwd) \
		&& export TMP_PRJ_DIR=/tmp$$PRJ_DIR \
		&& mkdir -p $$TMP_PRJ_DIR \
		&& cp -r $$PRJ_DIR $$TMP_PRJ_DIR/.. \
		&& cd $$TMP_PRJ_DIR \
		&& export GOPATH=/tmp$$GOPATH \
		&& dep ensure -v \
		&& diff $$PRJ_DIR/Gopkg.lock Gopkg.lock \
		&& diff -r $$PRJ_DIR/vendor vendor \
	'

.PHONY: test
test: test-unit test-api-compliance test-service-lifecycles

# Containerized unit tests-- requires docker-compose
.PHONY: test-unit
test-unit: check-docker-compose
	@# As of Go 1.9.0, testing ./... excludes tests on vendored code
	docker-compose run --rm test bash -c 'go test -tags unit ./...'

# Containerized service lifecycle tests-- requires docker-compose
.PHONY: test-service-lifecycles
test-service-lifecycles: check-docker-compose check-azure-env-vars
	@echo
	##############################################################################
	# WARNING! This creates services in Azure and will cost you real MONEY!      #
	# If run to completion, these tests clean up after themselves, but if        #
	# interrupted, you may need to perform some manual cleanup on your           #
	# subscription!                                                              #
	##############################################################################
	@echo
	docker-compose run \
		--rm \
		-e AZURE_SUBSCRIPTION_ID=$${AZURE_SUBSCRIPTION_ID} \
		-e AZURE_TENANT_ID=$${AZURE_TENANT_ID} \
		-e AZURE_CLIENT_ID=$${AZURE_CLIENT_ID} \
		-e AZURE_CLIENT_SECRET=$${AZURE_CLIENT_SECRET} \
		test \
		bash -c 'go test \
			-parallel 10 \
		  -timeout 60m \
			github.com/Azure/open-service-broker-azure/tests/lifecycle -v'


# Containerized API compliance check via osb-checker. Currently ignores exit code. 
.PHONY: test-api-compliance
test-api-compliance: check-docker-compose
	docker-compose build test-api-compliance-broker test-api-compliance
	-docker-compose run --rm test-api-compliance
	docker-compose kill test-api-compliance-broker
	docker-compose rm -f test-api-compliance-broker
	
.PHONY: lint
lint: check-docker-compose
	docker-compose run \
		--rm dev \
		bash -c 'gometalinter ./... \
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
			--deadline 120s \
			--vendor'

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

# (Re)Build the Docker image for the osba and run it
.PHONY: run
run: check-docker-compose check-azure-env-vars build
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
	docker build -t $(RC_IMAGE_NAME) .
	docker tag $(RC_IMAGE_NAME) $(RC_MUTABLE_IMAGE_NAME)

# Push the release candidate Docker images
.PHONY: docker-push-rc
docker-push-rc: check-docker docker-build
	docker push $(RC_IMAGE_NAME)
	docker push $(RC_MUTABLE_IMAGE_NAME)

# Push the release  / semver Docker images
.PHONY: docker-push-release
docker-push-release:
ifndef REL_VERSION
	$(error REL_VERSION is undefined)
endif
	docker pull $(RC_IMAGE_NAME)
	docker tag $(RC_IMAGE_NAME) $(REL_IMAGE_NAME)
	docker tag $(RC_IMAGE_NAME) $(REL_MUTABLE_IMAGE_NAME)

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

.PHONY: build-linux-broker-cli
build-linux-broker-cli: check-docker-compose
	docker-compose run --rm -e GOOS=linux -e GOARCH=amd64 dev \
		go build -o ${CONTRIB_BINARY_DIR}/${CLI_BINARY_NAME} \
		./contrib/cmd/cli

.PHONY: build-win-broker-cli
build-win-broker-cli: check-docker-compose
	docker-compose run --rm -e GOOS=windows -e GOARCH=amd64 dev \
		go build -o ${CONTRIB_BINARY_DIR}/${CLI_BINARY_NAME} \
		./contrib/cmd/cli
