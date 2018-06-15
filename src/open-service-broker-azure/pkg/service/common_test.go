package service

import (
	"fmt"

	"open-service-broker-azure/pkg/crypto/noop"
)

const fooValue = "bar"

var (
	testArbitraryObjectJSON = []byte(fmt.Sprintf(`{"foo":"%s"}`, fooValue))
	noopCodec               = noop.NewCodec()
)
