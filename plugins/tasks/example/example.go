package example

import (
	"github.com/influxdata/opsagent"
	"github.com/jaddqiu/opsagent/plugins/tasks"
)

// example.go

type Example struct {
	Ok bool
}

func (s *Example) Description() string {
	return "a demo plugin"
}

func (s *Example) SampleConfig() string {
	return `
  ## Indicate if everything is fine
  ok = true
`
}

func (s *Example) Gather(acc opsagent.Accumulator) error {
	if s.Ok {
		acc.AddFields("state", map[string]interface{}{"value": "pretty good"}, nil)
	} else {
		acc.AddFields("state", map[string]interface{}{"value": "not great"}, nil)
	}

	return nil
}

func init() {
	tasks.Add("example", func() opsagent.Task { return &Example{} })
}
