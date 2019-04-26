package client

import (
	"github.com/go-acme/lego/log"
)

// Check that every cert from the config is present in the file system
func ensureCertsFromConfig(storage *CertStorage, wantedDomains []string) func() {
	return func() {
		storedDomains, err := storage.LocallyStoredDomains()

		if err != nil {
			log.Warnf("Could not fetch stored certs: %v", err)
			return
		}

		// Loop over domains that shall be present
		for _, wantedDomain := range wantedDomains {

			if !stringInSlice(wantedDomain, storedDomains) {
				req := NewCertReq(wantedDomain, false)
				req.Submit()
			}
		}
	}
}

func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}
