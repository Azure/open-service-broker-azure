package service

import "fmt"

type ArbitraryType struct {
	Foo string `json:"foo"`
}

const fooValue = "bar"

var (
	testArbitraryObject = &ArbitraryType{
		Foo: fooValue,
	}
	testArbitraryObjectJSON = fmt.Sprintf(`{"foo":"%s"}`, fooValue)
)
