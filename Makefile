################################################################################
# Version details                                                              #
################################################################################

GIT_VERSION := $(shell git describe --always --abbrev=7 --dirty)

ifeq ($(REL_VERSION),)
	BROKER_VERSION := devel
else
	BROKER_VERSION := $(REL_VERSION)
endif

################################################################################
# Go build details                                                             #
################################################################################

BASE_PACKAGE_NAME := github.com/Azure/open-service-broker-azure

LDFLAGS := -w -X $(BASE_PACKAGE_NAME)/pkg/version.commit=$(GIT_VERSION) \
	-X $(BASE_PACKAGE_NAME)/pkg/version.version=$(BROKER_VERSION)

################################################################################
# Containerized development environment                                        #
################################################################################

DEV_IMAGE := quay.io/deis/lightweight-docker-go:v0.2.0

DOCKER_CMD_BASE := docker run \
	--rm \
	-e AZURE_SUBSCRIPTION_ID=$${AZURE_SUBSCRIPTION_ID} \
	-e AZURE_TENANT_ID=$${AZURE_TENANT_ID} \
	-e AZURE_CLIENT_ID=$${AZURE_CLIENT_ID} \
	-e AZURE_CLIENT_SECRET=$${AZURE_CLIENT_SECRET} \
	-e TEST_MODULES=$${TEST_MODULES} \
	-v $$(pwd):/go/src/$(BASE_PACKAGE_NAME) \
	-w /go/src/$(BASE_PACKAGE_NAME)

DOCKER_CMD := $(DOCKER_CMD_BASE) $(DEV_IMAGE)

DOCKER_CMD_INT := $(DOCKER_CMD_BASE) -it $(DEV_IMAGE)

################################################################################
# Docker images we build and publish                                           #
################################################################################

# This is left as 'azure-service-broker' because we don't yet have a docker repo
# for 'open-service-broker-azure'
#
# See https://github.com/Azure/open-service-broker-azure/issues/100
BASE_IMAGE_NAME := azure-service-broker

RC_IMAGE_NAME          := $(DOCKER_REPO)$(BASE_IMAGE_NAME):$(GIT_VERSION)
RC_MUTABLE_IMAGE_NAME  := $(DOCKER_REPO)$(BASE_IMAGE_NAME):canary

REL_IMAGE_NAME         := $(DOCKER_REPO)$(BASE_IMAGE_NAME):$(REL_VERSION)
REL_MUTABLE_IMAGE_NAME := $(DOCKER_REPO)$(BASE_IMAGE_NAME):latest

################################################################################
# Utility targets                                                              #
################################################################################

# Checks to ensure that AZURE_* environment variables needed to run the broker
# and its various test suites are set
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

# Allow developers to step into the containerized development environment--
# unconditionally requires docker
.PHONY: dev
dev:
	$(DOCKER_CMD_INT) bash

DEP_CMD := dep ensure -v

# Install/update dependencies
.PHONY: dep
dep:
ifdef SKIP_DOCKER
	$(DEP_CMD)
else
	$(DOCKER_CMD) $(DEP_CMD)
endif

################################################################################
# Tests                                                                        #
################################################################################

# Executes all tests
.PHONY: test
test: lint verify-vendored-code test-unit test-api-compliance \
	test-service-lifecycles

LINT_CMD := gometalinter ./... \
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
	--vendor

# Executes an extensive series of lint checks against broker code
.PHONY: lint
lint:
ifdef SKIP_DOCKER
	$(LINT_CMD)
else
	$(DOCKER_CMD) $(LINT_CMD)
endif

VERIFY_CMD := bash -c ' \
	export PRJ_DIR=$$(pwd) \
	&& export TMP_PRJ_DIR=/tmp$$PRJ_DIR \
	&& mkdir -p $$TMP_PRJ_DIR \
	&& cp -r $$PRJ_DIR $$TMP_PRJ_DIR/.. \
	&& cd $$TMP_PRJ_DIR \
	&& export GOPATH=/tmp$$GOPATH \
	&& dep ensure -v \
	&& diff $$PRJ_DIR/Gopkg.lock Gopkg.lock \
	&& diff -r $$PRJ_DIR/vendor vendor'

# Verifies there are no disrepancies between desired dependencies and the
# tracked, vendored dependencies
.PHONY: verify-vendored-code
verify-vendored-code:
ifdef SKIP_DOCKER
	$(VERIFY_CMD)
else
	$(DOCKER_CMD) $(VERIFY_CMD)
endif

# As of Go 1.9.0, testing ./... excludes tests on vendored code
UNIT_TEST_CMD := go test -tags unit ./...

# Executes unit tests
.PHONY: test-unit
test-unit:
ifdef SKIP_DOCKER
	$(UNIT_TEST_CMD)
else
	docker-compose run --rm test $(UNIT_TEST_CMD)
endif

LIFECYCLE_TEST_CMD := go test \
	-parallel 10 \
	-timeout 60m \
	$(BASE_PACKAGE_NAME)/tests/lifecycle -v

# Executes all or a subset of integration tests that test modules independently
# from the broker core/framework
.PHONY: test-service-lifecycles
test-service-lifecycles: check-azure-env-vars
	@echo
	##############################################################################
	# WARNING! This creates services in Azure and will cost you real MONEY!      #
	# If run to completion, these tests clean up after themselves, but if        #
	# interrupted, you may need to perform some manual cleanup on your           #
	# subscription!                                                              #
	##############################################################################
	@echo
ifdef SKIP_DOCKER
	$(LIFECYCLE_TEST_CMD)
else
	$(DOCKER_CMD) $(LIFECYCLE_TEST_CMD)
endif

# Evaluates broker compliance with the OSB specification-- unconditionally
# requires docker-compose
.PHONY: test-api-compliance
test-api-compliance:
	docker-compose build test-api-compliance-broker test-api-compliance
	-docker-compose run --rm test-api-compliance
	docker-compose kill test-api-compliance-broker
	docker-compose rm -f test-api-compliance-broker

################################################################################
# Misc                                                                         #
################################################################################

# Build the broker binary and docker image from code, then run it--
# unconditionally requires docker-compose
.PHONY: run
run: check-azure-env-vars
	@# Force the docker-compose "broker" service to be rebuilt-- this is separate
	@# from the build task used to produce a correctly tagged Docker image,
	@# although both builds are based on the same Dockerfile
	docker-compose build \
		--build-arg BASE_PACKAGE_NAME='$(BASE_PACKAGE_NAME)' \
		--build-arg LDFLAGS='$(LDFLAGS)' \
		broker
	docker-compose run \
		--rm -p 8080:8080 \
		-e AZURE_SUBSCRIPTION_ID=${AZURE_SUBSCRIPTION_ID} \
		-e AZURE_TENANT_ID=${AZURE_TENANT_ID} \
		-e AZURE_CLIENT_ID=${AZURE_CLIENT_ID} \
		-e AZURE_CLIENT_SECRET=${AZURE_CLIENT_SECRET} \
		broker

################################################################################
# Build / Publish                                                              #
################################################################################

# Build the broker binary and docker image
.PHONY: build
build:
	docker build \
		--build-arg BASE_PACKAGE_NAME='$(BASE_PACKAGE_NAME)' \
		--build-arg LDFLAGS='$(LDFLAGS)' \
		-t $(RC_IMAGE_NAME) \
		.
	docker tag $(RC_IMAGE_NAME) $(RC_MUTABLE_IMAGE_NAME)

# Push release candidate image
.PHONY: push-rc
push-rc: build
	docker push $(RC_IMAGE_NAME)
	docker push $(RC_MUTABLE_IMAGE_NAME)

# Push officially released, semantically versioned image
.PHONY: push-release
push-release:
ifndef REL_VERSION
	$(error REL_VERSION is undefined)
endif
	docker pull $(RC_IMAGE_NAME)
	docker tag $(RC_IMAGE_NAME) $(REL_IMAGE_NAME)
	docker tag $(RC_IMAGE_NAME) $(REL_MUTABLE_IMAGE_NAME)
	docker push $(REL_IMAGE_NAME)
	docker push $(REL_MUTABLE_IMAGE_NAME)
