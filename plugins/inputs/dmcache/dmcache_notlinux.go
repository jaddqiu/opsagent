// +build !linux

package dmcache

import (
	"github.com/jaddqiu/opsagent"
)

func (c *DMCache) Gather(acc opsagent.Accumulator) error {
	return nil
}

func dmSetupStatus() ([]string, error) {
	return []string{}, nil
}
