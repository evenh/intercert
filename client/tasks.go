package client

import (
	"context"
	"github.com/evenh/intercert/api"
	"github.com/go-acme/lego/v3/log"
	"io"
	"os"
)

func pingServer(client api.CertificateIssuerClient) func() {
	return func() {
		_, err := client.Ping(context.Background(), &api.PingRequest{Msg: "ping"})

		if err != nil {
			log.Warnf("Could not ping intercert host: %v", err)
		}
	}
}

// Listen for events from the server, indicating that a certificate has been renewed with ACME
func watchForEvents(domains []string, client api.CertificateIssuerClient) func() {
	return func() {
		// TODO: Automatic resubscribe when a stream dies for some reason (e.g. server goes down)
		renewalStream, err := client.OnCertificateRenewal(context.Background(), &api.CertificateRenewalNotificationRequest{
			DnsNames: domains,
		})

		if err != nil {
			log.Fatalf("Could not subscribe to renewal events: %v", err)
			os.Exit(1)
		}

		log.Infof("Listening for certificate renewal events")

		go func() {
			for {
				in, err := renewalStream.Recv()

				if err == io.EOF {
					break
				}

				if err != nil {
					log.Warnf("Got error while listening for renewal events: %v", err)
					break
				}

				log.Infof("Got notice from server that certificate for %s has been renewed. Queuing up re-fetch!", in.DnsName)
				job := NewCertReq(in.DnsName, true)
				job.Submit()
			}
		}()

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
