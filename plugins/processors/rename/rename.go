package rename

import (
	"github.com/jaddqiu/opsagent"
	"github.com/jaddqiu/opsagent/plugins/processors"
)

const sampleConfig = `
`

type Replace struct {
	Measurement string `toml:"measurement"`
	Tag         string `toml:"tag"`
	Field       string `toml:"field"`
	Dest        string `toml:"dest"`
}

type Rename struct {
	Replaces []Replace `toml:"replace"`
}

func (r *Rename) SampleConfig() string {
	return sampleConfig
}

func (r *Rename) Description() string {
	return "Rename measurements, tags, and fields that pass through this filter."
}

func (r *Rename) Apply(in ...opsagent.Metric) []opsagent.Metric {
	for _, point := range in {
		for _, replace := range r.Replaces {
			if replace.Dest == "" {
				continue
			}

			if replace.Measurement != "" {
				if value := point.Name(); value == replace.Measurement {
					point.SetName(replace.Dest)
				}
				continue
			}

			if replace.Tag != "" {
				if value, ok := point.GetTag(replace.Tag); ok {
					point.RemoveTag(replace.Tag)
					point.AddTag(replace.Dest, value)
				}
				continue
			}

			if replace.Field != "" {
				if value, ok := point.GetField(replace.Field); ok {
					point.RemoveField(replace.Field)
					point.AddField(replace.Dest, value)
				}
				continue
			}
		}
	}

	return in
}

func init() {
	processors.Add("rename", func() opsagent.Processor {
		return &Rename{}
	})
}
