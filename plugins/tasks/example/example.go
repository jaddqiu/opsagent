package example

import (
	_ "fmt"
    "log"
	"github.com/jaddqiu/opsagent"
	"github.com/jaddqiu/opsagent/plugins/tasks"
)

// example.go

// Example : task config for Example task
type Example struct {
	Ok bool
}

// Description : Description for Exmaple task
func (t *Example) Description() string {
	return "a demo plugin"
}

// SampleConfig : sample config for Exmaple task
func (t *Example) SampleConfig() string {
	return `
  ## Indicate if everything is fine
  ok = true
`
}

// Check : check envirionment for Example task
func (t *Example) Check() error {
	return nil
}

// Notify : notify execution result for Example task
func (t *Example) Notify() error {
	return nil
}

// Execute : task execution for Example task
func (t *Example) Execute() error {
	if t.Ok {
		log.Printf("Example task executed with ok set to true")
	} else {
		log.Printf("Example task executed with ok set to false")
	}

	return nil
}

func init() {
	tasks.Add("example", func() opsagent.Task { return &Example{} })
}
