package discard

import (
	"github.com/jaddqiu/opsagent"
	"github.com/jaddqiu/opsagent/plugins/outputs"
)

type Discard struct{}

func (d *Discard) Connect() error       { return nil }
func (d *Discard) Close() error         { return nil }
func (d *Discard) SampleConfig() string { return "" }
func (d *Discard) Description() string  { return "Send metrics to nowhere at all" }
func (d *Discard) Write(metrics []opsagent.Metric) error {
	return nil
}

func init() {
	outputs.Add("discard", func() opsagent.Output { return &Discard{} })
}
