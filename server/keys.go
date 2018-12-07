package server

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/xenolf/lego/registration"
	"io/ioutil"
	"strings"
)

const rsaKeySize = 2048

func CreatePrivateKey() (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, rsaKeySize)

	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func ReadPrivateKey(storageDirectory string, email string) (*rsa.PrivateKey, error) {


	data, err := ioutil.ReadFile(storageDirectory + "/keys/" + encodeEmail(email) + ".key")

	if err != nil {
		return nil, errors.New("could not find any existing private key for " + email)
	}

	block, _ := pem.Decode([]byte(data))

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)

	if err != nil {
		return nil, errors.New("could not parse private key")
	}

	return key, nil
}

func WritePrivateKey(storageDirectory string, email string, key *rsa.PrivateKey) error {
	if key == nil {
		return errors.New("missing private key - nothing to write to disk")
	}

	bytes := x509.MarshalPKCS1PrivateKey(key)
	privatePem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: bytes,
		},
	)

	err := ioutil.WriteFile(storageDirectory + "/keys/" + encodeEmail(email) + ".key", privatePem, 0600)

	if err != nil {
		log.Fatalf("Error writing private key: %v", err)
		return errors.New("could not write private key")
	}

	return nil
}

func ReadRegistration(storageDirectory string, email string) (*registration.Resource, error) {


	data, err := ioutil.ReadFile(storageDirectory + "/registration/" + encodeEmail(email))

	if err != nil {
		return nil, errors.New("could not find any existing registration for " + email)
	}

	var resource = registration.Resource{}

	err = json.Unmarshal(data, &resource)

	if err != nil {
		return nil, errors.New("could not unserialize registration data")
	}

	return &resource, nil
}

func WriteRegistration(storageDirectory string, email string, reg *registration.Resource) error {
	data, err := json.Marshal(reg)

	if err != nil {
		return errors.New("could not serialize registration")
	}

	err = ioutil.WriteFile(storageDirectory + "/registration/" + encodeEmail(email), data, 0600)

	if err != nil {
		log.Fatalf("Error writing registration: %v", err)
		return errors.New("could not write registration")
	}

	return nil
}

func encodeEmail(email string) string {
	h := sha256.New()
	h.Write([]byte(strings.ToLower(email)))
	return hex.EncodeToString(h.Sum(nil))
}
