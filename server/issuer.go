package server

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"fmt"
	"sort"
	"time"

	"github.com/go-acme/lego/certcrypto"

	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"strings"

	"github.com/evenh/intercert/api"
	"github.com/evenh/intercert/config"
	"github.com/go-acme/lego/log"
	"github.com/go-acme/lego/providers/dns"
	"github.com/mholt/certmagic"
	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// IssuerService issues certificates to clients
type IssuerService struct {
	client        *certmagic.Config
	whitelist     Whitelist
	renewalEvents chan string
}

// NewIssuerService constructs a new instance with a predefined config
func NewIssuerService(config *config.ServerConfig, userAgent string) *IssuerService {
	certmagic.UserAgent = userAgent
	
	issuer := new(IssuerService)
	issuer.renewalEvents = make(chan string)

	// Configure DNS provider by delegating to xenolf/lego factory
	dnsProvider, err := dns.NewDNSChallengeProviderByName(config.DNSProvider)

	if err != nil {
		log.Fatalf("Could not configure DNS provider: %v", err)
	}

	certmagicConfig := &certmagic.Config{
		CA:                      config.Directory,
		Email:                   config.Email,
		Agreed:                  config.Agree,
		DisableHTTPChallenge:    true,
		DisableTLSALPNChallenge: true,
		KeyType:                 certcrypto.RSA4096,
		MustStaple:              false,
		DNSProvider:             dnsProvider,
		Storage:                 createStorage(config.Storage),
		RenewDurationBefore:     config.RenewalThreshold,
		OnEvent: func(eventName string, payload interface{}) {
			// For now, we'll only care about certificates that are renewed
			if eventName == "acme_cert_renewed" {
				issuer.renewalEvents <- fmt.Sprintf("%s", payload)
			}
		},
	}

	// Construct the new certmagic instance
	magic := certmagic.New(certmagic.NewCache(certmagic.CacheOptions{
		RenewCheckInterval: config.ExpiryCheckAt,
		OCSPCheckInterval:  1 * time.Hour,
		GetConfigForCert: func(certificate certmagic.Certificate) (certmagic.Config, error) {
			return *certmagicConfig, nil
		},
	}), *certmagicConfig)

	issuer.client = magic

	// Create our whitelist
	whitelist := NewWhitelist(config.Domains)
	issuer.whitelist = whitelist

	log.Infof("Certificate issuer service configured - certificates will be renewed %v before expiry", config.RenewalThreshold)

	return issuer
}

// IssueCert issues a certificate for a valid request
func (s IssuerService) IssueCert(ctx context.Context, req *api.CertificateRequest) (*api.CertificateResponse, error) {
	// TODO: Validate auth in context
	logClient(ctx, "IssueCert("+req.DnsName+")")

	log.Infof("[%s] Received certificate request from client", req.DnsName)

	// Check whitelist
	err := s.whitelist.isDNSNameAllowed(req.DnsName)

	if err != nil {
		log.Warnf("[%s] Request rejected: %v", req.DnsName, err)
		return nil, err
	}

	// Hand over to certmagic
	err = s.client.Manage([]string{req.DnsName})

	if err != nil {
		log.Warnf("[%s] Failed to obtain certificate: %v", req.DnsName, err)
		return nil, err
	}

	cert, err := s.client.CacheManagedCertificate(req.DnsName)

	if err != nil {
		log.Warnf("[%s] Could not obtain certificate from ACME: %v", req.DnsName, err)
		return nil, errors.New("Could not obtain certificate from ACME")
	}

	// PEM-encode private key
	privateKey, err := pemEncodeKey(cert.PrivateKey)

	if err != nil {
		log.Warnf("[%s] Could not decode private key: %v", req.DnsName, err)
		return nil, errors.New("Could not decode private key")
	}

	// PEM-encode cert chain
	pemCert, err := pemEncodeCerts(cert.Certificate)

	if err != nil {
		log.Warnf("[%s] Could not PEM encode certificates: %v", req.DnsName, err)
		return nil, errors.New("Could not PEM encode certificates")
	}

	response := &api.CertificateResponse{
		Certificate: pemCert,
		PrivateKey:  string(privateKey),
		Names:       cert.Names,
	}

	log.Infof("[%s] Responding to client with certificate and private key", req.DnsName)

	return response, nil
}

// Ping is used to check that the service is alive
func (s IssuerService) Ping(ctx context.Context, req *api.PingRequest) (*api.PingResponse, error) {
	logClient(ctx, "Ping")
	// TODO: Auth for ping?
	return &api.PingResponse{Msg: "pong"}, nil
}

// Notifies the client about certificates that has been renewed server side
func (s IssuerService) OnCertificateRenewal(req *api.CertificateRenewalNotificationRequest, res api.CertificateIssuer_OnCertificateRenewalServer) error {
	logClient(res.Context(), "OnCertificateRenewal")
	names := req.DnsNames
	sort.Strings(names)

	for event := range s.renewalEvents {
		for _, name := range names {
			if strings.ToLower(name) == strings.ToLower(event) {
				err := res.Send(&api.RenewedCertificateEvent{DnsName: event})

				if err != nil {
					log.Warnf("[%s] Could not send renewal event to clients", name)
					return err
				}

				log.Infof("[%s] Notified clients that certificate has been renewed", name)
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

func logClient(ctx context.Context, operation string) {
	md, mdOK := metadata.FromIncomingContext(ctx)
	peerInfo, pOK := peer.FromContext(ctx)

	if mdOK && pOK {
		log.Infof("Call from %s - %s: %s", peerInfo.Addr, md["user-agent"], operation)
	}
}

func createStorage(dataDirectory string) certmagic.Storage {
	storage := &certmagic.FileStorage{Path: dataDirectory}
	log.Infof("Using directory %s for storage", storage)

	return storage
}
