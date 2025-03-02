package server

import (
	"net"
	"strconv"

	api "github.com/evenh/intercert/pkg/api/v1alpha1"
	"github.com/evenh/intercert/pkg/config"
	"github.com/go-acme/lego/v4/log"
	"google.golang.org/grpc"
)

// StartServer spawns a server instance given a server config
func StartServer(config *config.ServerConfig, userAgent string) {
	s := grpc.NewServer()

	issuerService := NewIssuerService(config, userAgent)

	api.RegisterCertificateIssuerServiceServer(s, issuerService)

	ln, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(config.Port))

	if err != nil {
		log.Fatal(err)
	}

	if err := s.Serve(ln); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
