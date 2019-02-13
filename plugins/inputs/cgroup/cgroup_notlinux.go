// +build !linux

package cgroup

import (
	"github.com/jaddqiu/opsagent"
)

func (g *CGroup) Gather(acc opsagent.Accumulator) error {
	return nil
}
