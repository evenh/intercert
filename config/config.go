package config

import (
	"strconv"
	"time"
)

// ServerConfig holds configuration that is required for running a server instance.
type ServerConfig struct {
	// Which port the server shall run on.
	Port int
	// Whether the user agrees to the Terms of Service
	// provided by the ACME vendor.
	Agree bool
	// The ACME directory to use.
	Directory string
	// The DNS provider that shall provide validation for
	// ACME DNS challenges.
	DNSProvider string
	// Whitelisted top level domains - e.g.
	// foo.com, bar.net.
	Domains []string
	// The email address of the server operator.
	Email string
	// The location on disk to save certificates and other data
	// that the server produces.
	Storage string
	// How often to check for expired certificates
	ExpiryCheckAt time.Duration
	// How early before expiry shall certificates be renewed?
	RenewalThreshold time.Duration
}

// ClientConfig holds configuration that is required for creating a client
type ClientConfig struct {
	// Which host shall the client connect to?
	Hostname string
	// Which port is the host listening on?
	Port int
	// The location on disk to save certificates and other data
	// that the client receives.
	Storage string
	// The domains to request certs for
	Domains []string
}

// GetDialAddr gets the formatted address to dial a new gRPC connection
func (c *ClientConfig) GetDialAddr() string {
	return c.Hostname + ":" + strconv.Itoa(c.Port)
}

// GetCertStorage returns the location on disk for where to store certificates
func (c *ClientConfig) GetCertStorage() string {
	return c.Storage + "/certs"
}
