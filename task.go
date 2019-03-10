package opsagent

type Task interface {
	// SampleConfig returns the default configuration of the Task
	SampleConfig() string

	// Description returns a one-sentence description on the Task
	Description() string

	// Execute task once.
	Execute() error

	// Notify the task execution result
	Notify() error
}
