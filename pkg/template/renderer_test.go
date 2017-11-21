package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type arbitraryType struct {
	Foo string
}

func TestRender(t *testing.T) {
	tpl := []byte(`{{ .Foo }}`)
	obj := arbitraryType{Foo: "bar"}
	out, err := Render(tpl, obj)
	assert.Nil(t, err)
	assert.Equal(t, obj.Foo, string(out))
}
