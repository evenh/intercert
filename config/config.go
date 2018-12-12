package config

// Holds configuration that is required for running a server instance.
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
	DnsProvider string
	// Whitelisted top level domains - e.g.
	// foo.com, bar.net.
	Domains []string
	// The email address of the server operator.
	Email string
	// The location on disk to save certificates and other data
	// that the server produces.
	Storage string
}
