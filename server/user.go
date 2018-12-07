package server

import (
	"crypto"
	log "github.com/sirupsen/logrus"
	"github.com/xenolf/lego/acme"
	"github.com/xenolf/lego/registration"
)

type AcmeUser struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func (u AcmeUser) GetEmail() string {
	return u.Email
}

func (u AcmeUser) GetRegistration() *registration.Resource {
	return u.Registration
}

func (u AcmeUser) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

func (u AcmeUser) LoadOrCreatePrivateKey(storageDirectory string) crypto.PrivateKey {
	if u.key == nil {
		key, err := ReadPrivateKey(storageDirectory, u.Email)

		if err != nil {
			// Generate new key
			log.Info("No existing private key found - creating a new one")
			key, err = CreatePrivateKey()

			if err != nil {
				panic(err)
			}

			// Save private key
			log.Info("Writing new private key")
			err := WritePrivateKey(storageDirectory, u.Email, key)

			if err != nil {
				panic(err)
			}

			log.Infof("Using the newly created private key for: %s", u.Email)
			return key
		} else {
			log.Infof("Loaded existing private key for: %s", u.Email)
			return key
		}
	}

	return u.key
}

func (u AcmeUser) LoadOrCreateRegistration(storageDirectory string, client *acme.Client) *registration.Resource {
	if u.Registration == nil {
		// Load existing one
		existingReg, err := ReadRegistration(storageDirectory, u.Email)

		if err != nil {
			log.Info("No existing registration found - registering with ACME")
			newReg, err := client.Registration.Register(true)

			if err != nil {
				log.Fatalf("Could not handle ACME registration: %v", err)
				return nil
			}

			// Save private key
			log.Info("Writing new registration")
			err = WriteRegistration(storageDirectory, u.Email, newReg)

			if err != nil {
				panic(err)
			}

			log.Infof("Using the newly created registration")
			return newReg
		}

		log.Infof("Loaded existing registration")
		return existingReg
	}

	return u.Registration
}
