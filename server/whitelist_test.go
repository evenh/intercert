package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleCase(t *testing.T) {
	w := NewWhitelist([]string{"github.com", "evenh.net"})

	err := w.isDNSNameAllowed("evenh.net")
	assert.NoError(t, err)

	err = w.isDNSNameAllowed("foo.evenh.net")
	assert.NoError(t, err)

	err = w.isDNSNameAllowed("bar.foo.evenh.net")
	assert.NoError(t, err)

	err = w.isDNSNameAllowed("*.evenh.net")
	assert.NoError(t, err)

	err = w.isDNSNameAllowed("github.com")
	assert.NoError(t, err)
}

func TestEmptyDomainList(t *testing.T) {
	w := NewWhitelist([]string{})

	err := w.isDNSNameAllowed("google.com")

	assert.NoError(t, err, "any domain should be allowed, including google.com")
}

func TestInvalidTopLevel(t *testing.T) {
	w := NewWhitelist([]string{"irrelevant"})

	err := w.isDNSNameAllowed("bar")

	assert.Error(t, err, "bar is not a valid top level domain")
	assert.EqualError(t, err, "Could not check whether bar is allowed")
}

func TestNoMatchInWhitelist(t *testing.T) {
	w := NewWhitelist([]string{"github.com", "evenh.net"})

	err := w.isDNSNameAllowed("foo.com")
	assert.Error(t, err)

	err = w.isDNSNameAllowed("foo.bar.net")
	assert.Error(t, err)

	err = w.isDNSNameAllowed("google.com")
	assert.Error(t, err)
}
