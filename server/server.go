package server

import (
	"net"
	"strconv"

	"github.com/evenh/intercert/api"
	"github.com/evenh/intercert/config"
	"github.com/go-acme/lego/log"
	"google.golang.org/grpc"
)

// StartServer spawns a server instance given a server config
func StartServer(config *config.ServerConfig, userAgent string) {
	s := grpc.NewServer()

	issuerService := NewIssuerService(config, userAgent)

	api.RegisterCertificateIssuerServer(s, issuerService)

	ln, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(config.Port))

	if err != nil {
		log.Fatal(err)
	}

	if err := s.Serve(ln); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
