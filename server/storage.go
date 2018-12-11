package server

import (
	"github.com/evenh/intercert/config"
	"github.com/mholt/certmagic"
	"github.com/xenolf/lego/log"
)

func IntercertCache(config *config.ServerConfig) *certmagic.Cache {
	var storage = createStorage(config.Storage);
	log.Infof("Using directory %s for storage", config.Storage)

	return createCache(storage)
}

func createCache(storage certmagic.Storage) *certmagic.Cache {
	return certmagic.NewCache(storage)
}

func createStorage(dataDirectory string) certmagic.Storage {
	return certmagic.FileStorage{Path: dataDirectory}
}
