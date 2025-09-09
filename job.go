package main

import "fmt"

type Job struct {
	ID      int    // Unique identifier for the job
	Payload string // Data to be processed
}

// String returns a string representation of the job
func (j Job) String() string {
	return fmt.Sprintf("Job{ID: %d, Payload: %s}", j.ID, j.Payload)
}
