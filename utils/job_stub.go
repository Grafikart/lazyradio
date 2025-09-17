//go:build !windows

package utils

import "os"

// Dummy JobObject so code compiles everywhere
type JobObject struct{}

func NewJobObject() (*JobObject, error) {
	return &JobObject{}, nil
}

func (j *JobObject) AddProcess(p *os.Process) error {
	return nil
}

func (j *JobObject) Close() error {
	return nil
}
