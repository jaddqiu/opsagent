package outputs

import (
	"github.com/jaddqiu/opsagent"
)

type Creator func() opsagent.Output

var Outputs = map[string]Creator{}

func Add(name string, creator Creator) {
	Outputs[name] = creator
}
