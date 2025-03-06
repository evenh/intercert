package client

import (
	"time"

	"go.uber.org/zap"
)

// From https://github.com/KyleBanks/go-kit/blob/master/job/job.go

// Job is a container for a repeating execution of a function.
//
// A Job executes it's function continuously with a predefined interval between
// executions. The Job only stops when the `Stop` function is called.
type Job struct {
	fn             func()
	interval       time.Duration
	runImmediately bool
	firstRun       bool
	logger         *zap.Logger

	stop    chan bool
	started bool
}

// Stop halts the execution of the Job's function.
func (j *Job) Stop() {
	j.logger.Info("stopping background job")
	j.stop <- true
	j.started = false
}

// start begins the execution of the Job's function.
//
// Note: It is unsafe to call `start` more than once!
func (j *Job) start() {
	j.started = true

	// Check if runImmediately is set on the first run
	if j.firstRun && j.runImmediately {
		j.fn()
	}
	j.firstRun = false

	if j.interval > 0*time.Second {
		go func() {
			for {
				// Sleep for the predetermined time.
				time.Sleep(j.interval)

				select {
				// Check for the 'stop' signal.
				case <-j.stop:
					return
				// Execute the function.
				default:
					j.fn()
				}
			}
		}()
	}
}

// Register schedules a function for execution, to be invoked repeated with a interval of
// the value of i.
//
// If the runImmediately parameter is true, the function will execute immediately. Otherwise,
// it will be invoked first after the duration of i.
func Register(f func(), name string, interval time.Duration, runImmediately bool, logger *zap.Logger) *Job {
	j := Job{
		fn:             f,
		interval:       interval,
		runImmediately: runImmediately,
		firstRun:       true,
		logger:         logger.With(zap.String("job", name)),

		stop: make(chan bool, 1),
	}

	logger.Info("configuring background job", zap.String("name", name), zap.Duration("interval", interval))

	j.start()
	return &j
}

// IsJobsStopped determines whether all jobs supplied is stopped or not
func IsJobsStopped(jobs []Job) bool {
	for _, job := range jobs {
		if job.started {
			return false
		}
	}

	return true
}
