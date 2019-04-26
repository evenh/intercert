package client

import (
	"context"
	"github.com/evenh/intercert/api"
	"github.com/go-acme/lego/log"
)

func pingServer(client api.CertificateIssuerClient) func() {
	return func() {
		_, err := client.Ping(context.Background(), &api.PingRequest{Msg: "ping"})

		if err != nil {
			log.Warnf("Could not ping intercert host: %v", err)
		}
	}
}

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
