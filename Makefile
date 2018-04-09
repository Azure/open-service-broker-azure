################################################################################
# Version details                                                              #
################################################################################

GIT_VERSION = $(shell git describe --always --abbrev=7 --dirty)

ifeq ($(REL_VERSION),)
	BROKER_VERSION := devel
else
	BROKER_VERSION := $(REL_VERSION)
endif

################################################################################
# Go build details                                                             #
################################################################################

BASE_PACKAGE_NAME := github.com/Azure/open-service-broker-azure

LDFLAGS = -w -X $(BASE_PACKAGE_NAME)/pkg/version.commit=$(GIT_VERSION) \
	-X $(BASE_PACKAGE_NAME)/pkg/version.version=$(BROKER_VERSION)

################################################################################
# Containerized development environment                                        #
################################################################################

DEV_IMAGE := quay.io/deis/lightweight-docker-go:v0.2.0

DOCKER_CMD_BASE := docker run \
	--rm \
	-e CGO_ENABLED=0 \
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
BASE_IMAGE_NAME        = azure-service-broker

RC_IMAGE_NAME          = $(DOCKER_REPO)$(BASE_IMAGE_NAME):$(GIT_VERSION)
RC_MUTABLE_IMAGE_NAME  = $(DOCKER_REPO)$(BASE_IMAGE_NAME):canary

REL_IMAGE_NAME         = $(DOCKER_REPO)$(BASE_IMAGE_NAME):$(REL_VERSION)
REL_MUTABLE_IMAGE_NAME = $(DOCKER_REPO)$(BASE_IMAGE_NAME):latest

################################################################################
# Utility targets                                                              #
################################################################################

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

# Deletes any existing osba binaries AND destroys any running Docker compose
# containers
.PHONY: clean
clean:
	rm -rf ${CONTRIB_BINARY_DIR}
ifndef SKIP_DOCKER
	docker-compose down --rmi local &> /dev/null
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

# Running the tests starts a containerized Redis dedicated to testing (if it
# isn't already running). It's left running afterwards (to speed up the next
# execution). It remains running unless explicitly shut down. This is a
# convenience task for stopping it.
.PHONY: stop-test-redis
stop-test-redis:
	docker-compose kill test-redis
	docker-compose rm -f test-redis

# Running the broker starts a containerized Redis dedicated to that purpose (if
# it isn't already running). It's left running afterwards (to speed up the next
# execution AND to retain state so broker recovery from incomplete async
# operations can be demonstrated). It remains running unless explicitly shut
# down. This is a convenience task for stopping it.
.PHONY: stop-broker-redis
stop-broker-redis:
	docker-compose kill broker-redis
	docker-compose rm -f broker-redis

################################################################################
# Tests                                                                        #
################################################################################

VERIFY_CMD := bash -c ' \
	PRJ_DIR=$$(pwd) \
	&& cp -r --parent -L $$PRJ_DIR /tmp \
	&& cd /tmp$$PRJ_DIR \
	&& GOPATH=/tmp$$GOPATH dep ensure -v \
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

.PHONY: test
test: test-unit test-api-compliance test-service-lifecycles

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

# Evaluates broker compliance with the OSB specification
.PHONY: test-api-compliance
test-api-compliance:
ifdef SKIP_DOCKER
	-/app/test.sh localhost 8088 60
else
	docker-compose build test-api-compliance-broker
	-docker-compose run --rm test-api-compliance
	docker-compose kill test-api-compliance-broker
	docker-compose rm -f test-api-compliance-broker
endif
	
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
	--deadline 240s \
	--vendor

# Executes an extensive series of lint checks against broker code
.PHONY: lint
lint:
ifdef SKIP_DOCKER
	$(LINT_CMD)
else
	$(DOCKER_CMD) $(LINT_CMD)
endif

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

# Rebuild and push officially released, semantically versioned image with
# semantically versioned binary
.PHONY: push-release
push-release:
ifndef REL_VERSION
	$(error REL_VERSION is undefined)
endif
	@# This pull is a verification that this commit has successfully cleared the
	@# master pipeline.
	docker pull $(RC_IMAGE_NAME)
	docker build \
		--build-arg BASE_PACKAGE_NAME='$(BASE_PACKAGE_NAME)' \
		--build-arg LDFLAGS='$(LDFLAGS)' \
		-t $(REL_IMAGE_NAME) \
		.
	docker tag $(REL_IMAGE_NAME) $(REL_MUTABLE_IMAGE_NAME)
	docker push $(REL_IMAGE_NAME)
	docker push $(REL_MUTABLE_IMAGE_NAME)

################################################################################
# Chart-Related Targets                                                        #
################################################################################

HELM_IMAGE := quay.io/deis/helm-chart-publishing-tools:v0.1.0

DOCKER_HELM_CMD := docker run \
	--rm \
	-e AZURE_STORAGE_CONNECTION_STRING=$${AZURE_STORAGE_CONNECTION_STRING} \
	-v $$(pwd):/go/src/$(BASE_PACKAGE_NAME) \
	-w /go/src/$(BASE_PACKAGE_NAME) \
	$(HELM_IMAGE)

LINT_CHART_CMD := helm lint contrib/k8s/charts/open-service-broker-azure \
	--set azure.tenantId=foo \
	--set azure.subscriptionId=foo \
	--set azure.clientId=foo \
	--set azure.clientSecret=foo

.PHONY: lint-chart
lint-chart:
ifdef SKIP_DOCKER
	$(LINT_CHART_CMD)
else
	$(DOCKER_HELM_CMD) $(LINT_CHART_CMD)
endif

PUBLISH_RC_CHART_CMD := bash -c ' \
	cd contrib/k8s/charts \
	&& rm -rf repo \
	&& mkdir repo \
	&& cd repo \
	&& sed -i "s/^version:.*/version: 0.0.1+$(GIT_VERSION)/g" ../open-service-broker-azure/Chart.yaml \
	&& sed -i "s/^appVersion:.*/appVersion: 0.0.1+$(GIT_VERSION)/g" ../open-service-broker-azure/Chart.yaml \
	&& sed -i "s/^  tag:.*/  tag: $(GIT_VERSION)/g" ../open-service-broker-azure/values.yaml \
	&& helm dep build ../open-service-broker-azure \
	&& helm package ../open-service-broker-azure \
	&& az storage blob upload \
		-c azure-rc \
		--file open-service-broker-azure-0.0.1+$(GIT_VERSION).tgz \
		--name open-service-broker-azure-0.0.1+$(GIT_VERSION).tgz \
	&& az storage container lease acquire -c azure-rc --lease-duration 60 \
	&& helm repo index --url https://kubernetescharts.blob.core.windows.net/azure-rc . \
	&& az storage blob upload \
		-c azure-rc \
		--file index.yaml \
		--name index.yaml'

PUBLISH_RELEASE_CHART_CMD := bash -c ' \
	SIMPLE_REL_VERSION=$$(echo $(REL_VERSION) | cut -c 2-) \
	&& cd contrib/k8s/charts \
	&& rm -rf repo \
	&& mkdir repo \
	&& cd repo \
	&& sed -i "s/^version:.*/version: $${SIMPLE_REL_VERSION}/g" ../open-service-broker-azure/Chart.yaml \
	&& sed -i "s/^appVersion:.*/appVersion: $${SIMPLE_REL_VERSION}/g" ../open-service-broker-azure/Chart.yaml \
	&& sed -i "s/^  tag:.*/  tag: $(REL_VERSION)/g" ../open-service-broker-azure/values.yaml \
	&& helm dep build ../open-service-broker-azure \
	&& helm package ../open-service-broker-azure \
	&& az storage blob upload \
		-c azure \
		--file open-service-broker-azure-$${SIMPLE_REL_VERSION}.tgz \
		--name open-service-broker-azure-$${SIMPLE_REL_VERSION}.tgz \
	&& az storage container lease acquire -c azure --lease-duration 60 \
	&& az storage blob download \
		-c azure \
		--name index.yaml \
		--file index.yaml \
	&& helm repo index --url https://kubernetescharts.blob.core.windows.net/azure --merge index.yaml . \
	&& az storage blob upload \
		-c azure \
		--file index.yaml \
		--name index.yaml'

.PHONY: publish-rc-chart
publish-rc-chart:
ifndef AZURE_STORAGE_CONNECTION_STRING
	$(error AZURE_STORAGE_CONNECTION_STRING is not defined)
endif
ifdef SKIP_DOCKER
	$(PUBLISH_RC_CHART_CMD)
else
	$(DOCKER_HELM_CMD) $(PUBLISH_RC_CHART_CMD)
endif

.PHONY: publish-release-chart
publish-release-chart:
ifndef REL_VERSION
	$(error REL_VERSION is undefined)
endif
ifndef AZURE_STORAGE_CONNECTION_STRING
	$(error AZURE_STORAGE_CONNECTION_STRING is not defined)
endif
ifdef SKIP_DOCKER
	$(PUBLISH_RELEASE_CHART_CMD)
else
	$(DOCKER_HELM_CMD) $(PUBLISH_RELEASE_CHART_CMD)
endif

################################################################################
# contrib/                                                                     #
################################################################################

CONTRIB_BINARY_DIR := contrib/bin
CLI_BINARY_NAME := broker-cli

BUILD_CMD := go build -o $(CONTRIB_BINARY_DIR)/$(CLI_BINARY_NAME) ./cmd/cli

.PHONY: build-broker-cli
build-broker-cli:
ifdef SKIP_DOCKER
	$(BUILD_CMD)
else
	docker run \
	--rm \
	-e GOOS=$(GOOS) \
	-e GOARCH=$(GOARCH) \
	-v $$(pwd)/..:/go/src/$(BASE_PACKAGE_NAME) \
	-w /go/src/$(BASE_PACKAGE_NAME)/contrib \
	$(DEV_IMAGE) \
	$(BUILD_CMD)
endif

################################################################################
# PCF Tile Targets                                                             #
################################################################################

TILE_GENERATOR_IMAGE := cfplatformeng/tile-generator:v11.0.4

DOCKER_TILE_CMD := docker run \
	--rm \
	-v $$(pwd):/workspace \
	-w /workspace \
	$(TILE_GENERATOR_IMAGE)

GENERATE_TILE_CMD := bash -c ' \
	rm -rf \
		contrib/cf/pcf-tile/product \
		contrib/cf/pcf-tile/release \
		contrib/cf/pcf-tile/resources/open-service-broker-azure.zip \
		contrib/cf/pcf-tile/tile-history.yml \
	&& zip -r \
		contrib/cf/pcf-tile/resources/open-service-broker-azure.zip \
		cmd pkg vendor \
	&& cd contrib/cf/pcf-tile \
	&& tile build $$(echo $(REL_VERSION) | cut -c 2-)'

.PHONY: generate-pcf-tile
generate-pcf-tile:
ifdef SKIP_DOCKER
	$(GENERATE_TILE_CMD)
else
	$(DOCKER_TILE_CMD) $(GENERATE_TILE_CMD)
endif
