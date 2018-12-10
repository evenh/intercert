package server

import (
	"context"
	"github.com/evenh/intercert/api"
	"github.com/evenh/intercert/config"
	log "github.com/sirupsen/logrus"
	"github.com/xenolf/lego/acme"
	"github.com/xenolf/lego/providers/dns"
)

type IssuerService struct {
	user AcmeUser
	client *acme.Client
}

func NewIssuerService(config *config.ServerConfig) *IssuerService {
	issuer := new(IssuerService)

	// User based on email and private key
	user := AcmeUser {Email: config.Email}
	key := user.LoadOrCreatePrivateKey(config.Storage)
	user.key = key

	// Create client
	log.Infof("Using directory server: %s", config.Directory)
	client, err := acme.NewClient(config.Directory, user, acme.RSA2048)

	if err != nil {
		log.Fatalf("Could not construct new ACME client: %v", err)
		log.Exit(1)
	}

	// Configure DNS provider
	provider, err:= dns.NewDNSChallengeProviderByName(config.DnsProvider)

	if err != nil {
		log.Fatalf("Could not configure DNS provider: %v", err)
	}

	// Remove other challenges
	client.ExcludeChallenges([]acme.Challenge{ acme.HTTP01, acme.TLSALPN01 })

	// Only use DNS challenge
	_ = client.SetChallengeProvider(acme.DNS01, provider)

	// Handle registration
	reg := user.LoadOrCreateRegistration(config.Storage, client)

	if reg == nil {
		log.Fatalf("Failed to obtain registration")
		log.Exit(1)
	}

	user.Registration = reg
	issuer.user = user
	issuer.client = client

	return issuer
}

func (s IssuerService) IssueCert(ctx context.Context, req *api.CertificateRequest) (*api.CertificateResponse, error) {
	// TODO: Validate auth in context

	certificates, err := s.client.ObtainCertificate([]string{ req.DnsName }, true, nil, false)

	if err != nil {
		log.Warnf("Failed to obtain certificate: %v", err)
		return nil, err
	}

	log.Infof("Received payload: %#v", certificates)

	response := &api.CertificateResponse{CertificatePayload: certificates.Certificate}

	return response, nil
}

