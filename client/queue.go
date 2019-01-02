package client

import (
	"context"
	"fmt"
	"github.com/evenh/intercert/api"
	"github.com/xenolf/lego/log"
)

var JobQueue = make(chan CertReq, 10)

type CertReq struct {
	DnsName  string
	Renewal  bool
	Attempts int
}

func NewCertReq(dnsName string, renewal bool) CertReq {
	return CertReq{
		DnsName:  dnsName,
		Renewal:  renewal,
		Attempts: 0,
	}
}

func (r *CertReq) Submit() {
	var operation string

	if r.Renewal {
		operation = "renewal"
	} else {
		operation = "initial certificates"
	}

	log.Infof("Submitting '%s' to job queue for %s", r.DnsName, operation)
	JobQueue <- *r
}

func CertWorker(id int, client api.CertificateIssuerClient, storage *CertStorage) {
	var prefix = fmt.Sprintf("[worker-%d]", id)
	for job := range JobQueue {

		if job.Attempts >= 5 {
			log.Warnf("%s Attempts for %s exceed 5, discarding job", prefix, job.DnsName)
			return
		}

		log.Infof("%s Processing %s", prefix, job.DnsName)

		res, err := client.IssueCert(context.TODO(), &api.CertificateRequest{DnsName: job.DnsName})

		if err != nil {
			log.Warnf("%s Encountered error while requesting certs for %s: %v", prefix, job.DnsName, err)
			job.Attempts++
			log.Warnf("%s Increasing attempt counter to %d and re-queuing...", prefix, job.Attempts)
			JobQueue <- job
		} else {
			log.Infof("%s Successfully fetched certs for %s", prefix, job.DnsName)
			err := storage.Store(job.DnsName, res)

			if err != nil {
				log.Warnf("Could not write cert data for %s: %v", job.DnsName, err)
			}
			// TODO: Store certs
		}
	}
}
