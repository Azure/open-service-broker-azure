################################################################################
# Go build details                                                             #
################################################################################

BINARY_DIR := bin

BINARY_NAME := broker-cli

BASE_PACKAGE_NAME := github.com/Azure/open-service-broker-azure

################################################################################
# Containerized development environment                                        #
################################################################################

DEV_IMAGE := quay.io/deis/lightweight-docker-go:v0.2.0

DOCKER_CMD := docker run \
	--rm \
	-e GOOS=$(GOOS) \
	-e GOARCH=$(GOARCH) \
	-v $$(pwd)/..:/go/src/$(BASE_PACKAGE_NAME) \
	-w /go/src/$(BASE_PACKAGE_NAME)/contrib \
	$(DEV_IMAGE)

################################################################################
# Build                                                                        #
################################################################################

.PHONY: clean
clean:
	rm -rf ${BINARY_DIR}

BUILD_CMD := go build -o $(BINARY_DIR)/$(BINARY_NAME) ./cmd/cli

.PHONY: build-broker-cli
build-broker-cli:
ifdef SKIP_DOCKER
	$(BUILD_CMD)
else
	$(DOCKER_CMD) $(BUILD_CMD)
endif
