package printer

import (
	"fmt"

	"github.com/jaddqiu/opsagent"
	"github.com/jaddqiu/opsagent/plugins/processors"
	"github.com/jaddqiu/opsagent/plugins/serializers"
	"github.com/jaddqiu/opsagent/plugins/serializers/influx"
)

type Printer struct {
	serializer serializers.Serializer
}

var sampleConfig = `
`

func (p *Printer) SampleConfig() string {
	return sampleConfig
}

func (p *Printer) Description() string {
	return "Print all metrics that pass through this filter."
}

func (p *Printer) Apply(in ...telegraf.Metric) []telegraf.Metric {
	for _, metric := range in {
		octets, err := p.serializer.Serialize(metric)
		if err != nil {
			continue
		}
		fmt.Printf("%s", octets)
	}
	return in
}

func init() {
	processors.Add("printer", func() telegraf.Processor {
		return &Printer{
			serializer: influx.NewSerializer(),
		}
	})
}
