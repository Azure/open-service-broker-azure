package service

import (
	"fmt"
)

const fooValue = "bar"

var (
	testArbitraryObjectJSON = []byte(fmt.Sprintf(`{"foo":"%s"}`, fooValue))
)
