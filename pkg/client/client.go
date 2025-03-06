package client

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	api "github.com/evenh/intercert/pkg/api/v1alpha1"
	"github.com/evenh/intercert/pkg/config"
	"github.com/go-acme/lego/v4/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
)

var logger *zap.Logger

// StartClient starts a client instance with a client config and a user agent
func StartClient(config *config.ClientConfig, userAgent string, l *zap.Logger) {
	logger = l

	logger.Info("initializing client")

	// Check if the config is valid
	err := validateConfig(config)

	if err != nil {
		logger.Fatal("invalid configuration", zap.Error(err))
	}

	logger.Info("configuring connection to server for gRPC operations", zap.String("address", config.GetDialAddr()))

	// Configure connection
	// TODO: Use TLS for encrypting connection (with configurable certs and a sane default)
	conn, err := grpc.NewClient(config.GetDialAddr(),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.DefaultConfig,
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUserAgent(userAgent+";"),
	)

	if err != nil {
		logger.Fatal("could not configure connection to host: %v", zap.Error(err))
	}

	defer conn.Close()

	// Create client from connection
	client := api.NewCertificateIssuerServiceClient(conn)

	// Test connection
	_, err = client.Ping(context.Background(), &api.PingRequest{Msg: "ping"})

	if err != nil {
		logger.Fatal("could not test connection", zap.String("address", config.GetDialAddr()), zap.Error(err))
	}

	logger.Info("successfully pinged intercert host")
	logger.Info("client running - ready to work!")

	// Create storage for saving/loading certs
	certStorage := NewCertStorage(config.GetCertStorage(), logger.Named("storage"))

	// Configure cert workers
	for w := 1; w <= runtime.NumCPU(); w++ {
		go CertWorker(w, client, certStorage, logger.Named("queue.worker"))
	}

	// Set up scheduled tasks
	tasks := configureTasks(client, config, certStorage)

	// Handle termination
	configureTermination(tasks)

	// Block forever
	select {}
}

func configureTermination(backgroundJobs []Job) {
	var gracefulStop = make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-gracefulStop
		log.Infof("Caught signal: %+v", sig)
		log.Infof("Waiting for background tasks to exit...")

		// Tell all background jobs to exit
		for i := range backgroundJobs {
			backgroundJobs[i].Stop()
		}

		// Wait for completion of background jobs
		for !IsJobsStopped(backgroundJobs) {
		}

		os.Exit(0)
	}()
}

func validateConfig(c *config.ClientConfig) error {
	if len(c.Domains) < 1 {
		return errors.New("no domains specified")
	}

	if len(c.Hostname) == 0 {
		return errors.New("hostname was empty")
	}

	return nil
}

func configureTasks(client api.CertificateIssuerServiceClient, config *config.ClientConfig, storage *CertStorage) []Job {
	var tasks []Job

	pinger := *Register(pingServer(client), "Ping intercert host", 10*time.Minute, false, logger.Named("tasks.ping"))
	renewalHandler := *Register(watchForEvents(config.Domains, client), "Watch for certificate renewal events", 0*time.Second, true, logger.Named("tasks.renewal"))
	desiredCheck := *Register(ensureCertsFromConfig(storage, config.Domains), "Ensure configured domains is present", 1*time.Hour, true, logger.Named("tasks.initial-certs"))

	tasks = append(tasks, pinger, renewalHandler, desiredCheck)

	return tasks
}
