// +build !linux

package wireless

import (
	"log"

	"github.com/jaddqiu/opsagent"
	"github.com/jaddqiu/opsagent/plugins/inputs"
)

func (w *Wireless) Gather(acc telegraf.Accumulator) error {
	return nil
}

func init() {
	inputs.Add("wireless", func() telegraf.Input {
		log.Print("W! [inputs.wireless] Current platform is not supported")
		return &Wireless{}
	})
}
