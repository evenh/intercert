package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleCase(t *testing.T) {
	w := Whitelist{[]string{"foo.com"}}

	err := w.isDnsNameAllowed("bar.foo.com")

	assert.Nil(t, err, "bar.foo.com should be a match of foo.com")
}
