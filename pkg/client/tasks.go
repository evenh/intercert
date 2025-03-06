package client

import (
	"context"
	"io"
	"slices"

	api "github.com/evenh/intercert/pkg/api/v1alpha1"
	"go.uber.org/zap"
)

func pingServer(client api.CertificateIssuerServiceClient) func() {
	return func() {
		_, err := client.Ping(context.Background(), &api.PingRequest{Msg: "ping"})

		if err != nil {
			logger.Warn("could not ping intercert host", zap.Error(err))
		}
	}
}

// Listen for events from the server, indicating that a certificate has been renewed with ACME
func watchForEvents(domains []string, client api.CertificateIssuerServiceClient) func() {
	return func() {
		// TODO: Automatic resubscribe when a stream dies for some reason (e.g. server goes down)
		renewalStream, err := client.OnCertificateRenewal(context.Background(), &api.OnCertificateRenewalRequest{
			DnsNames: domains,
		})

		if err != nil {
			logger.Fatal("could not watch for certificate renewal", zap.Error(err))
		}

		logger.Info("listening for certificate renewal events")

		go func() {
			for {
				in, err := renewalStream.Recv()

				if err == io.EOF {
					break
				}

				if err != nil {
					logger.Warn("got error while listening for renewal events", zap.Error(err))
					break
				}

				logger.Info("got notice from server that certificate has been renewed – queuing up a re-fetch job", zap.String("dns_name", in.DnsName))
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
			logger.Warn("could not fetch stored certs", zap.Error(err))
			return
		}

		// Loop over domains that shall be present
		for _, wantedDomain := range wantedDomains {

			if !slices.Contains(storedDomains, wantedDomain) {
				req := NewCertReq(wantedDomain, false)
				req.Submit()
			}
		}
	}
}
