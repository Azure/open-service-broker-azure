SHORT_NAME := azure-service-broker

BINARY_DEST_DIR := bin

.PHONY: test
test:
	go test $$(glide nv)

.PHONY: build
build:
	go build -o ${BINARY_DEST_DIR}/${SHORT_NAME} ./pkg

.PHONY: run
run: build
	${BINARY_DEST_DIR}/${SHORT_NAME}

