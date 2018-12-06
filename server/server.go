package server

import (
	"github.com/evenh/intercert/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"strconv"
)

func StartServer(port int) {
	s := grpc.NewServer()

	issuerService := IssuerService{}

	api.RegisterCertificateIssuerServer(s, issuerService)

	ln, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(port))

	if err != nil {
		log.Fatal(err)
	}

	// TODO: Remove after testing? Used for grpcui
	reflection.Register(s)

	if err := s.Serve(ln); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
