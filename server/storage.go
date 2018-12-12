package server

import (
	"github.com/evenh/intercert/config"
	"github.com/mholt/certmagic"
	"github.com/xenolf/lego/log"
)

// Creates a new cache, backed by a file system storage. This is more
// or less the same as provided by certmagic, but customized to save to
// directory provided in the configuration.
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
