package server

import (
	"errors"
	"strings"

	"github.com/bobesa/go-domain-util/domainutil"
	"go.uber.org/zap"
)

// whitelist holds valid domains that certificates can be
// issued under.
type whitelist struct {
	// The top level domains to allow
	domains []string
	logger  *zap.Logger
}

// newWhitelist creates a new instance of the whitelist struct
func newWhitelist(domains []string, logger *zap.Logger) *whitelist {
	lowercasedDomains := make([]string, len(domains))

	for i, v := range domains {
		lowercasedDomains[i] = strings.ToLower(v)
	}

	if len(lowercasedDomains) == 0 {
		logger.Warn("no domains in whitelist - every domain is allowed")
	} else {
		logger.Info("loaded whitelist", zap.Strings("domains", lowercasedDomains))
	}

	return &whitelist{
		domains: lowercasedDomains,
		logger:  logger,
	}
}

// Checks whether a DNS name (e.g. foo.bar.com) is allowed. If no domains is configured,
// every DNS name will be allowed.
func (w *whitelist) isDNSNameAllowed(dnsName string) error {
	// Empty whitelists allows everything
	if len(w.domains) == 0 {
		return nil
	}

	topLevel := domainutil.Domain(dnsName)

	if topLevel == "" {
		return errors.New("could not check whether " + dnsName + " is allowed")
	}

	for _, domain := range w.domains {
		if domain == topLevel {
			return nil
		}
	}

	return errors.New(topLevel + " did not match any entries in the whitelist")
}
