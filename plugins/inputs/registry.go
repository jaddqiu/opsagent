package inputs

import "github.com/jaddqiu/opsagent"

type Creator func() opsagent.Input

var Inputs = map[string]Creator{}

func Add(name string, creator Creator) {
	Inputs[name] = creator
}
