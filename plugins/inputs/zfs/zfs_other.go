// +build !linux,!freebsd

package zfs

import (
	"github.com/jaddqiu/opsagent"
	"github.com/jaddqiu/opsagent/plugins/inputs"
)

func (z *Zfs) Gather(acc telegraf.Accumulator) error {
	return nil
}

func init() {
	inputs.Add("zfs", func() telegraf.Input {
		return &Zfs{}
	})
}
