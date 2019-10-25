package client

import (
	"os"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/evenh/intercert/api"
	"github.com/go-acme/lego/v3/log"
	"github.com/pkg/errors"
)

const (
	certificate = iota
	privateKey
)

// CertStorage represents where the certificates is stored at the client side
type CertStorage struct {
	storageDirectory string
}

// NewCertStorage constructs an instance of the CertStorage struct, with validation
func NewCertStorage(storageDirectory string) *CertStorage {
	if _, err := os.Stat(storageDirectory); os.IsNotExist(err) {
		err = os.MkdirAll(storageDirectory, 0777)

		if err != nil {
			log.Warnf("Could not create directory for certs: %v", err)
			os.Exit(1)
		}
	}

	return &CertStorage{
		storageDirectory: storageDirectory,
	}
}

// LocallyStoredDomains lists all domain names that has TLS certificates stored
func (c *CertStorage) LocallyStoredDomains() ([]string, error) {
	existingCerts, err := c.readCerts()

	if err != nil {
		return nil, err
	}

	sort.Strings(existingCerts)

	return existingCerts, nil
}

// Store certificates for a given API response and domain
func (c *CertStorage) Store(domain string, response *api.CertificateResponse) error {
	// Find disk locations
	certFile, _ := c.absoluteFileName(domain, certificate)
	keyFile, _ := c.absoluteFileName(domain, privateKey)

	locationContent := map[string]string{
		certFile: response.Certificate,
		keyFile:  response.PrivateKey,
	}

	log.Infof("Attempting to write cert data for %s", domain)

	for location, content := range locationContent {
		fh, err := os.OpenFile(location, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)

		if err != nil {
			return err
		}

		_, err = fh.WriteString(content)

		if err != nil {
			return err
		}

		err = fh.Sync()

		if err != nil {
			return err
		}

		err = fh.Close()

		if err != nil {
			return err
		}
	}

	log.Infof("Successfully wrote cert data for %s", domain)
	return nil
}

func (c *CertStorage) absoluteFileName(domain string, keyType int) (string, error) {
	base := c.storageDirectory + "/" + strings.ToLower(domain)

	switch keyType {
	case privateKey:
		return base + ".key", nil
	case certificate:
		return base + ".cer", nil
	default:
		return "", errors.New("Unknown key type: " + strconv.Itoa(keyType))
	}
}

func (c *CertStorage) readCerts() ([]string, error) {
	var files []string
	f, err := os.Open(c.storageDirectory)

	if err != nil {
		return files, err
	}

	fileInfo, err := f.Readdir(-1)
	_ = f.Close()

	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		trimmed := strings.TrimSuffix(file.Name(), path.Ext(file.Name()))
		files = append(files, trimmed)
	}

	return unique(files), nil
}

func unique(input []string) []string {
	u := make([]string, 0, len(input))
	m := make(map[string]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}
