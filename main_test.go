package main

import (
	"sync"
	"testing"
	"time"
)

func TestJobProcessing(t *testing.T) {
	// Test configuration
	const (
		maxWorkers   = 2
		jobQueueSize = 5
		numJobs      = 4
	)

	// Create dispatcher
	dispatcher := NewDispatcher(maxWorkers, jobQueueSize)

	// Start dispatcher
	err := dispatcher.Start()
	if err != nil {
		t.Fatalf("Failed to start dispatcher: %v", err)
	}

	// Override the processJob function to count processed jobs
	// We'll use a different approach: submit jobs and wait for completion
	jobs := make([]Job, numJobs)
	for i := 0; i < numJobs; i++ {
		jobs[i] = Job{
			ID:      i + 1,
			Payload: "test-payload",
		}
	}

	// Submit all jobs
	for _, job := range jobs {
		err := dispatcher.SubmitJob(job)
		if err != nil {
			t.Errorf("Failed to submit job: %v", err)
		}
	}

	// Wait for jobs to be processed
	time.Sleep(1 * time.Second)

	// Stop dispatcher
	dispatcher.Stop()

	// Verify that we have the expected number of workers
	if len(dispatcher.WorkerPool) != maxWorkers {
		t.Errorf("Expected %d workers, got %d", maxWorkers, len(dispatcher.WorkerPool))
	}

	// Verify that the job queue was created with the correct size
	if cap(dispatcher.JobQueue) != jobQueueSize {
		t.Errorf("Expected job queue size %d, got %d", jobQueueSize, cap(dispatcher.JobQueue))
	}
}

func TestDispatcherInvalidWorkers(t *testing.T) {
	// Test with invalid worker count
	dispatcher := NewDispatcher(0, 10)
	err := dispatcher.Start()
	if err == nil {
		t.Error("Expected error when starting dispatcher with 0 workers")
	}

	dispatcher2 := NewDispatcher(-1, 10)
	err2 := dispatcher2.Start()
	if err2 == nil {
		t.Error("Expected error when starting dispatcher with negative workers")
	}
}

func TestJobQueueFull(t *testing.T) {
	// Test with very small queue to trigger full queue scenario
	const (
		maxWorkers   = 1
		jobQueueSize = 1
	)

	dispatcher := NewDispatcher(maxWorkers, jobQueueSize)
	err := dispatcher.Start()
	if err != nil {
		t.Fatalf("Failed to start dispatcher: %v", err)
	}

	// Fill the queue
	job1 := Job{ID: 1, Payload: "job1"}
	err = dispatcher.SubmitJob(job1)
	if err != nil {
		t.Errorf("Failed to submit first job: %v", err)
	}

	// Try to add another job immediately (queue should be full)
	job2 := Job{ID: 2, Payload: "job2"}
	err = dispatcher.SubmitJob(job2)
	if err == nil {
		t.Error("Expected error when submitting to full queue")
	}

	// Clean up
	time.Sleep(200 * time.Millisecond) // Let job process
	dispatcher.Stop()
}

func TestJobString(t *testing.T) {
	job := Job{ID: 123, Payload: "test-data"}
	expected := "Job{ID: 123, Payload: test-data}"
	if job.String() != expected {
		t.Errorf("Expected %s, got %s", expected, job.String())
	}
}

func TestWorkerCreation(t *testing.T) {
	var wg sync.WaitGroup
	worker := NewWorker(1, &wg)

	if worker.ID != 1 {
		t.Errorf("Expected worker ID 1, got %d", worker.ID)
	}

	if worker.JobChannel == nil {
		t.Error("Worker JobChannel should not be nil")
	}

	if worker.QuitChannel == nil {
		t.Error("Worker QuitChannel should not be nil")
	}
}
