package client

import (
	"context"
	"time"

	api "github.com/evenh/intercert/pkg/api/v1alpha1"
	"go.uber.org/zap"
)

const (
	maxAttempts     = 5
	issuanceTimeout = 15 * time.Minute
)

// jobQueue holds all active jobs
var jobQueue = make(chan certReq, 10)

// certReq represents a certificate request
type certReq struct {
	DNSName  string
	Renewal  bool
	Attempts int
}

// NewCertReq constructs an instance of the certReq struct
func NewCertReq(dnsName string, renewal bool) certReq {
	return certReq{
		DNSName:  dnsName,
		Renewal:  renewal,
		Attempts: 0,
	}
}

// Submit a job to the queue for processing
func (r *certReq) Submit() {
	logger.Named("queue").Info("submitting job", zap.String("dns_name", r.DNSName), zap.Bool("renewal", r.Renewal))
	jobQueue <- *r
}

// CertWorker is the actual worker routine that eats away from the queue
func CertWorker(id int, client api.CertificateIssuerServiceClient, storage *CertStorage, logger *zap.Logger) {
	for job := range jobQueue {
		jobLog := logger.With(
			zap.Int("worker_id", id),
			zap.String("dns_name", job.DNSName),
			zap.Bool("renewal", job.Renewal),
		)

		if job.Attempts >= maxAttempts {
			jobLog.Warn("max attempts reached, discarding job", zap.Int("attempts", job.Attempts))
			return
		}

		jobLog.Info("processing job")

		ctx, cancel := context.WithTimeout(context.Background(), issuanceTimeout)
		res, err := client.IssueCert(ctx, &api.IssueCertRequest{DnsName: job.DNSName})
		cancel()

		if err != nil {
			jobLog.Warn("encountered error while requesting cert", zap.Int("attempts", job.Attempts), zap.Error(err))
			job.Attempts++
			jobLog.Warn("increasing attempt counter to and re-queuing...", zap.Int("attempts", job.Attempts))
			jobQueue <- job
		} else {
			jobLog.Info("successfully fetched certs")
			err := storage.Store(job.DNSName, res)

			if err != nil {
				jobLog.Warn("could not write cert data – job and data discarded", zap.Error(err))
			}
		}
	}
}
