package client

import (
	"context"
	"errors"
	"github.com/evenh/intercert/api"
	"github.com/evenh/intercert/config"
	"github.com/xenolf/lego/log"
	"google.golang.org/grpc"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
)

func StartClient(config *config.ClientConfig) {
	log.Infof("Initializing client")

	// Check if the config is valid
	err := validateConfig(config)

	if err != nil {
		log.Warnf("Configuration error: %v", err)
		os.Exit(1)
	}

	log.Infof("Configuring connection to %s for gRPC operations", config.GetDialAddr())

	// Configure connection
	conn, err := grpc.Dial(config.GetDialAddr(), grpc.WithInsecure()) // TODO: Not run insecure

	if err != nil {
		log.Warnf("Could not configure connection to host: %v", err)
		os.Exit(1)
	}

	defer conn.Close()

	// Create client from connection
	client := api.NewCertificateIssuerClient(conn)

	// Test connection
	_, err = client.Ping(context.Background(), &api.PingRequest{Msg: "ping"})

	if err != nil {
		log.Warnf("Could not test connection to %s: %v", config.GetDialAddr(), err)
		os.Exit(1)
	}

	log.Infof("Successfully pinged intercert host")

	// Wait group for graceful shutdown
	wg := sync.WaitGroup{}

	// Create storage for saving/loading certs
	certStorage := NewCertStorage(config.GetCertStorage())

	// Configure cert workers
	for w := 1; w <= runtime.NumCPU(); w++ {
		go CertWorker(w, client, certStorage)
	}

	// Set up scheduled tasks
	tasks := configureTasks(config, &wg, certStorage)

	// Handle termination
	configureTermination(&wg, tasks)

	log.Infof("Client running - ready to work!")

	// Start process by scanning for desired certs
	ensureCertsFromConfig(certStorage, config.Domains)()

	// Block forever
	select {}
}

func configureTermination(wg *sync.WaitGroup, backgroundChannels []chan bool) {
	var gracefulStop = make(chan os.Signal)

	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	go func() {
		sig := <-gracefulStop
		log.Infof("Caught signal: %+v", sig)
		log.Infof("Waiting for background tasks to exit...")

		// Tell all background jobs to exit
		for i := range backgroundChannels {
			backgroundChannels[i] <- true
		}

		// Wait for completion of background jobs
		wg.Wait()

		os.Exit(0)
	}()
}

func validateConfig(c *config.ClientConfig) error {
	if len(c.Domains) <= 1 {
		return errors.New("no domains specified")
	}

	if len(c.Hostname) == 0 {
		return errors.New("hostname was empty")
	}

	if c.RenewalThreshold > (24*time.Hour)*30 {
		return errors.New("renewal threshold can't exceed 30 days")
	}

	return nil
}

func configureTasks(config *config.ClientConfig, wg *sync.WaitGroup, storage *CertStorage) []chan bool {
	var tasks []chan bool

	expiryCheck := schedule(findExpiredCerts(config.RenewalThreshold), "Expired Certs watcher", config.ExpiryCheckAt, wg)
	desiredCheck := schedule(ensureCertsFromConfig(storage, config.Domains), "Ensure configured domains is present", 8*time.Hour, wg)

	tasks = append(tasks, expiryCheck, desiredCheck)
	return tasks
}
