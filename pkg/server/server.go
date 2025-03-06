package server

import (
	"net"
	"strconv"

	api "github.com/evenh/intercert/pkg/api/v1alpha1"
	"github.com/evenh/intercert/pkg/config"
	"github.com/go-acme/lego/v4/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// StartServer spawns a server instance given a server config
func StartServer(config *config.ServerConfig, userAgent string, logger *zap.Logger) {
	s := grpc.NewServer()

	issuerService := NewIssuerService(config, userAgent, logger)

	api.RegisterCertificateIssuerServiceServer(s, issuerService)

	ln, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(config.Port))

	if err != nil {
		log.Fatal(err)
	}

	if err := s.Serve(ln); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
