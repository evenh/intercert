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

	issuer.user = user

	// Create client
	clientConfig := acme.NewConfig(issuer.user)
	clientConfig.CADirURL = config.Directory

	client, err := acme.NewClient(clientConfig)

	if err != nil {
		log.Fatal("Could not construct new ACME client: %v", err)
	} else {
		issuer.client = client
	}

	return issuer
}

func (s IssuerService) IssueCert(context.Context, *api.CertificateRequest) (*api.CertificateResponse, error) {
	// TODO: Validate auth in context
	log.Println("Got request for issuing")
	return nil, errors.New("dummy error")
}

