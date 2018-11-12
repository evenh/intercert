package cmd

// A ServerConfig represents configuration for launching a server instance.
type ServerConfig struct {
	Port        int
	Agree       bool
	Directory   string
	DnsProvider string
	Domains     []string
	Email       string
	Storage     string
}
