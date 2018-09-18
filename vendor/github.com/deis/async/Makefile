################################################################################
# Go build details                                                             #
################################################################################

BASE_PACKAGE_NAME := github.com/deis/async

################################################################################
# Containerized development environment                                        #
################################################################################

DEV_IMAGE := quay.io/deis/lightweight-docker-go:v0.2.0

DOCKER_CMD_BASE := docker run \
	--rm \
	-v $$(pwd):/go/src/$(BASE_PACKAGE_NAME) \
	-w /go/src/$(BASE_PACKAGE_NAME)

DOCKER_CMD := $(DOCKER_CMD_BASE) $(DEV_IMAGE)

DOCKER_CMD_INT := $(DOCKER_CMD_BASE) -it $(DEV_IMAGE)

################################################################################
# Utility targets                                                              #
################################################################################

# Destroys any running Docker compose containers
.PHONY: clean
clean:
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

# Running the tests starts a containerized Redis (if it isn't already running).
# It's left running afterwards (to speed up the next execution). It remains
# running unless explicitly shut down. This is a convenience task for stopping
# it.s
.PHONY: stop-test-redis
stop-test-redis:
	docker-compose kill test-redis
	docker-compose rm -f test-redis

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

# As of Go 1.9.0, testing ./... excludes tests on vendored code
TEST_CMD := go test ./...

# Executes tests
.PHONY: test
test:
ifdef SKIP_DOCKER
	$(TEST_CMD)
else
	docker-compose run --rm test $(TEST_CMD)
endif

LINT_CMD := gometalinter ./... \
	--concurrency=1 \
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
