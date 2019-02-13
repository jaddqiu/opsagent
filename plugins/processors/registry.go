package processors

import "github.com/jaddqiu/opsagent"

type Creator func() opsagent.Processor

var Processors = map[string]Creator{}

func Add(name string, creator Creator) {
	Processors[name] = creator
}
