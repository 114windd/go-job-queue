package main

import (
	"fmt"
	"sync"
	"time"
)

// Worker represents a worker that processes jobs
type Worker struct {
	ID          int             // Unique identifier for the worker
	JobChannel  chan Job        // Channel to receive jobs
	QuitChannel chan bool       // Channel to signal worker to quit
	wg          *sync.WaitGroup // WaitGroup to track worker completion
}

// NewWorker creates and returns a new Worker
func NewWorker(id int, wg *sync.WaitGroup) *Worker {
	return &Worker{
		ID:          id,
		JobChannel:  make(chan Job),
		QuitChannel: make(chan bool),
		wg:          wg,
	}
}

// Start begins the worker's job processing loop
func (w *Worker) Start() {
	// Signal that this goroutine has started
	w.wg.Add(1)

	go func() {
		// Ensure we signal completion when the goroutine ends
		defer w.wg.Done()

		fmt.Printf("Worker %d starting\n", w.ID)

		for {
			select {
			case job := <-w.JobChannel:
				// Process the received job
				w.processJob(job)
			case <-w.QuitChannel:
				// Quit signal received, exit the worker
				fmt.Printf("Worker %d stopping\n", w.ID)
				return
			}
		}
	}()
}

// Stop signals the worker to quit
func (w *Worker) Stop() {
	w.QuitChannel <- true
}

// processJob simulates job processing
func (w *Worker) processJob(job Job) {
	fmt.Printf("Worker %d processing %s\n", w.ID, job.String())

	// Simulate work by sleeping
	time.Sleep(time.Millisecond * 100)

	fmt.Printf("Worker %d completed %s\n", w.ID, job.String())
}
