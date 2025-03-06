package dnsprovider

import (
	"errors"

	"github.com/libdns/libdns"
)

var availableProviders = map[string]providerConfig{}

type DnsProvider interface {
	libdns.RecordGetter
	libdns.RecordAppender
	libdns.RecordSetter
	libdns.RecordDeleter
}

type providerConfig struct {
	name     string
	provider func() (DnsProvider, error)
}

func registerProvider(p providerConfig) {
	if len(p.name) == 0 {
		panic("failed registering provider: empty name")
	}

	if p.provider == nil {
		panic("failed registering provider: nil function")
	}

	availableProviders[p.name] = p
}

// TODO: caching/sync.Once ?
func GetProvider(name string) (DnsProvider, error) {
	p, ok := availableProviders[name]
	if !ok {
		return nil, errors.New("unknown provider: " + name)
	}

	return p.provider()
}
