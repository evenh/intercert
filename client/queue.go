package client

import (
	"context"
	"fmt"

	"github.com/evenh/intercert/api"
	"github.com/go-acme/lego/log"
)

// JobQueue holds all active jobs
var JobQueue = make(chan CertReq, 10)

// CertReq represents a certificate request
type CertReq struct {
	DNSName  string
	Renewal  bool
	Attempts int
}

// NewCertReq constructs an instance of the CertReq struct
func NewCertReq(dnsName string, renewal bool) CertReq {
	return CertReq{
		DNSName:  dnsName,
		Renewal:  renewal,
		Attempts: 0,
	}
}

// Submit a job to the queue for processing
func (r *CertReq) Submit() {
	var operation string

	if r.Renewal {
		operation = "renewal"
	} else {
		operation = "initial certificates"
	}

	log.Infof("Submitting '%s' to job queue for %s", r.DNSName, operation)
	JobQueue <- *r
}

// CertWorker is the actual worker routine that eats away from the queue
func CertWorker(id int, client api.CertificateIssuerClient, storage *CertStorage) {
	var prefix = fmt.Sprintf("[worker-%d]", id)
	for job := range JobQueue {

		if job.Attempts >= 5 {
			log.Warnf("%s Attempts for %s exceed 5, discarding job", prefix, job.DNSName)
			return
		}

		log.Infof("%s Processing %s", prefix, job.DNSName)

		res, err := client.IssueCert(context.TODO(), &api.CertificateRequest{DnsName: job.DNSName})

		if err != nil {
			log.Warnf("%s Encountered error while requesting certs for %s: %v", prefix, job.DNSName, err)
			job.Attempts++
			log.Warnf("%s Increasing attempt counter to %d and re-queuing...", prefix, job.Attempts)
			JobQueue <- job
		} else {
			log.Infof("%s Successfully fetched certs for %s", prefix, job.DNSName)
			err := storage.Store(job.DNSName, res)

			if err != nil {
				log.Warnf("Could not write cert data for %s: %v", job.DNSName, err)
			}
		}
	}
}
