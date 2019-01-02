package client

import (
	"github.com/xenolf/lego/log"
	"sync"
	"time"
)

func schedule(fn func(), name string, delay time.Duration, wg *sync.WaitGroup) chan bool {
	log.Infof("Configuring background task '%s' every %v", name, delay)
	stop := make(chan bool)

	go func() {
		for {
			select {
			case <-time.After(delay):
				fn()
			case <-stop:
				log.Infof("Stopping background task: %s", name)
				// Decrease the WaitGroup by 1
				wg.Done()
				return
			}
		}
	}()

	return stop
}
