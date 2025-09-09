package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	// Configuration
	const (
		maxWorkers   = 3
		jobQueueSize = 10
		numJobs      = 8
	)

	// Create dispatcher
	dispatcher := NewDispatcher(maxWorkers, jobQueueSize)

	// Start the dispatcher
	if err := dispatcher.Start(); err != nil {
		log.Fatalf("Failed to start dispatcher: %v", err)
	}

	// Submit jobs
	fmt.Printf("\nSubmitting %d jobs...\n", numJobs)
	for i := 1; i <= numJobs; i++ {
		job := Job{
			ID:      i,
			Payload: fmt.Sprintf("Task-%d", i),
		}

		if err := dispatcher.SubmitJob(job); err != nil {
			log.Printf("Failed to submit job: %v", err)
		}
	}

	// Let jobs process for a while
	fmt.Println("\nProcessing jobs...")
	time.Sleep(2 * time.Second)

	// Stop the dispatcher
	dispatcher.Stop()

}
