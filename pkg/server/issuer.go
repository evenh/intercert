package server

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/evenh/intercert/pkg/config"
	"github.com/evenh/intercert/pkg/dnsprovider"
	"go.uber.org/zap"

	"github.com/caddyserver/certmagic"
	api "github.com/evenh/intercert/pkg/api/v1alpha1"
	"github.com/go-acme/lego/v4/log"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

var (
	renewalEvents = make(chan string)
)

// IssuerService issues certificates to clients
type IssuerService struct {
	api.UnimplementedCertificateIssuerServiceServer
	client        *certmagic.Config
	renewalEvents *chan string
	whitelist     *whitelist
	logger        *zap.Logger
}

// NewIssuerService constructs a new instance with a predefined config
func NewIssuerService(config *config.ServerConfig, userAgent string, logger *zap.Logger) *IssuerService {
	certmagic.UserAgent = userAgent
	issuerLogger := logger.Named("issuer")
	cmLogger := issuerLogger.Named("certmagic")

	dnsProvider, err := dnsprovider.GetProvider(config.DNSProvider)
	storage := createLocalStorage(config.Storage)

	if err != nil {
		issuerLogger.Fatal("Could not configure DNS provider", zap.Error(err))
	}

	certmagicConfig := &certmagic.Config{
		Logger:     cmLogger,
		MustStaple: false,
		Storage:    storage,
		OnEvent: func(ctx context.Context, eventName string, data map[string]any) error {
			// For now, we'll only care about certificates that are renewed
			if eventName == "cert_obtaining" && data["renewal"] == true {
				renewalEvents <- fmt.Sprintf("%s", data["name"])
			}

			return nil
		},
	}

	// Must be constructed outside the config in order to have the mutex set
	acmeIssuer := certmagic.NewACMEIssuer(certmagicConfig, certmagic.ACMEIssuer{
		CA:                      config.Directory,
		Email:                   config.Email,
		Agreed:                  config.Agree,
		DisableHTTPChallenge:    true,
		DisableTLSALPNChallenge: true,
		Logger:                  cmLogger.Named("acme"),
		DNS01Solver: &certmagic.DNS01Solver{
			DNSManager: certmagic.DNSManager{
				DNSProvider: dnsProvider,
				Logger:      cmLogger.Named("dnsprovider"),
				TTL:         60,
				//PropagationDelay:   0,
				//PropagationTimeout: 0,
				//Resolvers:          nil,
				//OverrideDomain:     "",
			},
		},
	})
	certmagicConfig.Issuers = append(certmagicConfig.Issuers, acmeIssuer)

	// Construct the new certmagic instance
	magic := certmagic.New(certmagic.NewCache(certmagic.CacheOptions{
		RenewCheckInterval: config.ExpiryCheckAt,
		Logger:             cmLogger,
		GetConfigForCert: func(certificate certmagic.Certificate) (*certmagic.Config, error) {
			return certmagicConfig, nil
		},
	}), *certmagicConfig)

	issuer := &IssuerService{
		renewalEvents: &renewalEvents,
		client:        magic,
		logger:        issuerLogger,
		whitelist:     newWhitelist(config.Domains, logger.Named("whitelist")),
	}

	issuerLogger.Info("certificate issuer service configured",
		zap.String("acme_directory", config.Directory),
		zap.String("storage_path", config.Storage),
	)

	return issuer
}

// IssueCert issues a certificate for a valid request
func (s IssuerService) IssueCert(ctx context.Context, req *api.IssueCertRequest) (*api.IssueCertResponse, error) {
	// TODO: Validate auth in context
	logClient(ctx, "IssueCert("+req.DnsName+")")

	reqLog := s.logger.With(zap.String("dns_name", req.DnsName))

	reqLog.Info("received certificate request")

	err := s.whitelist.isDNSNameAllowed(req.DnsName)
	if err != nil {
		reqLog.Warn("dns name not allowed", zap.Error(err))
		return nil, err
	}

	// Hand over to certmagic
	// TODO: Consider ManageAsync
	err = s.client.ManageSync(ctx, []string{req.DnsName})

	if err != nil {
		reqLog.Warn("failed to obtain certificate", zap.Error(err))
		return nil, err
	}

	cert, err := s.client.CacheManagedCertificate(ctx, req.DnsName)

	if err != nil {
		reqLog.Warn("could not obtain certificate from ACME", zap.Error(err))
		return nil, errors.New("could not obtain certificate from ACME")
	}

	// PEM-encode private key
	privateKey, err := pemEncodeKey(cert.PrivateKey)

	if err != nil {
		reqLog.Warn("could not decode private key", zap.Error(err))
		return nil, errors.New("could not decode private key")
	}

	// PEM-encode cert chain
	pemCert, err := pemEncodeCerts(cert.Certificate)

	if err != nil {
		reqLog.Warn("could not PEM encode certificates", zap.Error(err))
		return nil, errors.New("could not PEM encode certificates")
	}

	response := &api.IssueCertResponse{
		Certificate: pemCert,
		PrivateKey:  string(privateKey),
		Names:       cert.Names,
	}

	reqLog.Info("responding to client with certificate and private key")

	return response, nil
}

// Ping is used to check that the service is alive
func (s IssuerService) Ping(ctx context.Context, req *api.PingRequest) (*api.PingResponse, error) {
	logClient(ctx, "Ping")
	// TODO: Auth for ping?
	return &api.PingResponse{Msg: "pong"}, nil
}

// OnCertificateRenewal notifies the client about certificates that has been renewed server side
func (s IssuerService) OnCertificateRenewal(req *api.OnCertificateRenewalRequest, res api.CertificateIssuerService_OnCertificateRenewalServer) error {
	logClient(res.Context(), "OnCertificateRenewal")
	names := req.DnsNames
	sort.Strings(names)

	for event := range *s.renewalEvents {
		for _, name := range names {
			if strings.ToLower(name) == strings.ToLower(event) {
				err := res.Send(&api.OnCertificateRenewalResponse{DnsName: event})

				if err != nil {
					s.logger.Warn("could not send renewal event to clients", zap.String("dns_name", name), zap.Error(err))
					return err
				}

				s.logger.Info("notified clients that certificate has been renewed", zap.String("dns_name", name))
			}
		}
	}

	return nil
}

// from certmagic: encodePrivateKey marshals a EC or RSA private key into a PEM-encoded array of bytes.
func pemEncodeKey(key crypto.PrivateKey) ([]byte, error) {
	var pemType string
	var privKeyBytes []byte

	switch key := key.(type) {
	case *ecdsa.PrivateKey:
		var err error
		pemType = "EC"
		privKeyBytes, err = x509.MarshalECPrivateKey(key)
		if err != nil {
			return nil, err
		}
	case *rsa.PrivateKey:
		pemType = "RSA"
		privKeyBytes = x509.MarshalPKCS1PrivateKey(key)
	}

	privatePemKey := pem.Block{Type: pemType + " PRIVATE KEY", Bytes: privKeyBytes}

	return pem.EncodeToMemory(&privatePemKey), nil
}

// create a string containing the PEM encoded certificate chain
func pemEncodeCerts(cert tls.Certificate) (string, error) {
	var certificates []string

	for _, blob := range cert.Certificate {
		pemBlock := pem.EncodeToMemory(&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: blob,
		})

		certificates = append(certificates, string(pemBlock))
	}

	return strings.Join(certificates, ""), nil
}

// TODO: Middleware
func logClient(ctx context.Context, operation string) {
	md, mdOK := metadata.FromIncomingContext(ctx)
	peerInfo, pOK := peer.FromContext(ctx)

	if mdOK && pOK {
		log.Infof("Call from %s - %s: %s", peerInfo.Addr, md["user-agent"], operation)
	}
}

func createLocalStorage(dataDirectory string) certmagic.Storage {
	storage := &certmagic.FileStorage{Path: dataDirectory}

	return storage
}
