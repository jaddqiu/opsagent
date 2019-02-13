package aggregators

import "github.com/jaddqiu/opsagent"

type Creator func() opsagent.Aggregator

var Aggregators = map[string]Creator{}

func Add(name string, creator Creator) {
	Aggregators[name] = creator
}
