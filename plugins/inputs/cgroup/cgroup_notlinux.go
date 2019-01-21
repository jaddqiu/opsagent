// +build !linux

package cgroup

import (
	"github.com/jaddqiu/opsagent"
)

func (g *CGroup) Gather(acc telegraf.Accumulator) error {
	return nil
}
