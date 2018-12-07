package server

import (
	"context"
	"errors"
	"github.com/evenh/intercert/api"
	"github.com/evenh/intercert/config"
	log "github.com/sirupsen/logrus"
	"github.com/xenolf/lego/acme"
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
	clientConfig := acme.NewConfig(user)
	clientConfig.CADirURL = config.Directory
	client, err := acme.NewClient(clientConfig)

	if err != nil {
		log.Fatalf("Could not construct new ACME client: %v", err)
		log.Exit(1)
	}

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

func (s IssuerService) IssueCert(context.Context, *api.CertificateRequest) (*api.CertificateResponse, error) {
	// TODO: Validate auth in context
	log.Println("Got request for issuing")
	return nil, errors.New("dummy error")
}

