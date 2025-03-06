package dnsprovider

import (
	"errors"
	"os"

	"github.com/libdns/digitalocean"
)

func init() {
	do := providerConfig{
		name: "digitalocean",
		provider: func() (DnsProvider, error) {
			t, found := os.LookupEnv("DO_AUTH_TOKEN")
			if !found {
				return nil, errors.New("DO_AUTH_TOKEN environment variable not found")
			}

			p := &digitalocean.Provider{
				APIToken: t,
			}

			return p, nil
		},
	}

	registerProvider(do)
}
