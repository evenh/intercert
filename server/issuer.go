package server

import (
	"context"
	"github.com/evenh/intercert/api"
	"github.com/evenh/intercert/config"
	"github.com/mholt/certmagic"
	"github.com/xenolf/lego/certcrypto"
	"github.com/xenolf/lego/log"
	"github.com/xenolf/lego/providers/dns"
)

type IssuerService struct {
	client *certmagic.Config
}

func NewIssuerService(config *config.ServerConfig) *IssuerService {
	issuer := new(IssuerService)

	// Configure DNS provider by delegating to xenolf/lego factory
	dnsProvider, err:= dns.NewDNSChallengeProviderByName(config.DnsProvider)

	if err != nil {
		log.Fatalf("Could not configure DNS provider: %v", err)
	}

	// Construct the new certmagic instance
	magic := certmagic.NewWithCache(IntercertCache(config), certmagic.Config{
		CA:     config.Directory,
		Email:  config.Email,
		Agreed: config.Agree,
		DisableHTTPChallenge: true,
		DisableTLSALPNChallenge: true,
		KeyType: certcrypto.RSA4096,
		MustStaple: false,
		DNSProvider: dnsProvider,
	})

	issuer.client = magic

	return issuer
}

func (s IssuerService) IssueCert(ctx context.Context, req *api.CertificateRequest) (*api.CertificateResponse, error) {
	// TODO: Validate auth in context

	log.Infof("[%s] Received certificate request from client", req.DnsName)

	err := s.client.Manage([]string{ req.DnsName })

	if err != nil {
		log.Warnf("Failed to obtain certificate: %v", err)
		return nil, err
	}

	cert, err := s.client.CacheManagedCertificate(req.DnsName)

	log.Infof("[%s] Received payload: %#v", req.DnsName, cert)

	response := &api.CertificateResponse{CertificatePayload: []byte{} }

	return response, nil
}
