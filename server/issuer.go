package server

import (
	"context"
	"errors"
	"github.com/evenh/intercert/api"
	"log"
)

type IssuerService struct{}

func (IssuerService) IssueCert(context.Context, *api.CertificateRequest) (*api.CertificateResponse, error) {
	log.Println("Got request for issuing")
	return nil, errors.New("dummy error")
}
