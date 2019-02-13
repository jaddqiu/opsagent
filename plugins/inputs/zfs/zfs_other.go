// +build !linux,!freebsd

package zfs

import (
	"github.com/jaddqiu/opsagent"
	"github.com/jaddqiu/opsagent/plugins/inputs"
)

func (z *Zfs) Gather(acc opsagent.Accumulator) error {
	return nil
}

func init() {
	inputs.Add("zfs", func() opsagent.Input {
		return &Zfs{}
	})
}
