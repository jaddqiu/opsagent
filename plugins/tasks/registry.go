package tasks

import "github.com/jaddqiu/opsagent"

type Creator func() opsagent.Task

var Tasks = map[string]Creator{}

func Add(name string, creator Creator) {
	Tasks[name] = creator
}
