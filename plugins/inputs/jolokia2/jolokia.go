package jolokia2

import (
	"github.com/jaddqiu/opsagent"
	"github.com/jaddqiu/opsagent/plugins/inputs"
)

func init() {
	inputs.Add("jolokia2_agent", func() opsagent.Input {
		return &JolokiaAgent{
			Metrics:               []MetricConfig{},
			DefaultFieldSeparator: ".",
		}
	})
	inputs.Add("jolokia2_proxy", func() opsagent.Input {
		return &JolokiaProxy{
			Metrics:               []MetricConfig{},
			DefaultFieldSeparator: ".",
		}
	})
}
