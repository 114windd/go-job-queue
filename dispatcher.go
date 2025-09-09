package main

import (
	"fmt"
	"sync"
)

// Dispatcher manages a pool of workers and distributes jobs
type Dispatcher struct {
	JobQueue   chan Job       // Channel to receive incoming jobs
	WorkerPool []*Worker      // Pool of available workers
	MaxWorkers int            // Maximum number of workers
	wg         sync.WaitGroup // WaitGroup to track worker completion
}

// NewDispatcher creates and returns a new Dispatcher
func NewDispatcher(maxWorkers int, jobQueueSize int) *Dispatcher {
	// Handle invalid parameters to prevent panics
	if maxWorkers < 0 {
		maxWorkers = 0
	}
	if jobQueueSize < 0 {
		jobQueueSize = 0
	}

	return &Dispatcher{
		JobQueue:   make(chan Job, jobQueueSize),
		WorkerPool: make([]*Worker, 0, maxWorkers),
		MaxWorkers: maxWorkers,
	}
}

// Start initializes the worker pool and begins job distribution
func (d *Dispatcher) Start() error {
	if d.MaxWorkers <= 0 {
		return fmt.Errorf("maxWorkers must be greater than 0")
	}

	fmt.Printf("Starting dispatcher with %d workers\n", d.MaxWorkers)

	// Create and start workers
	for i := 1; i <= d.MaxWorkers; i++ {
		worker := NewWorker(i, &d.wg)
		d.WorkerPool = append(d.WorkerPool, worker)
		worker.Start()
	}

	// Start the job distribution goroutine
	go d.dispatch()

	return nil
}

// dispatch distributes jobs from the job queue to available workers
func (d *Dispatcher) dispatch() {
	workerIndex := 0

	for job := range d.JobQueue {
		// Round-robin job distribution to workers
		if len(d.WorkerPool) > 0 {
			worker := d.WorkerPool[workerIndex%len(d.WorkerPool)]
			workerIndex++

			// Send job to the selected worker (non-blocking)
			go func(w *Worker, j Job) {
				w.JobChannel <- j
			}(worker, job)
		}
	}
}

// SubmitJob adds a job to the job queue
func (d *Dispatcher) SubmitJob(job Job) error {
	select {
	case d.JobQueue <- job:
		fmt.Printf("Submitted %s to job queue\n", job.String())
		return nil
	default:
		return fmt.Errorf("job queue is full, cannot submit %s", job.String())
	}
}

// Stop gracefully shuts down all workers
func (d *Dispatcher) Stop() {
	fmt.Println("Stopping dispatcher...")

	// Stop all workers
	for _, worker := range d.WorkerPool {
		worker.Stop()
	}

	// Wait for all workers to finish
	d.wg.Wait()

	fmt.Println("All workers stopped")
}

// Wait blocks until all workers have completed their current jobs
func (d *Dispatcher) Wait() {
	d.wg.Wait()
}
