package all

import (
	_ "github.com/jaddqiu/opsagent/plugins/processors/converter"
	_ "github.com/jaddqiu/opsagent/plugins/processors/enum"
	_ "github.com/jaddqiu/opsagent/plugins/processors/override"
	_ "github.com/jaddqiu/opsagent/plugins/processors/parser"
	_ "github.com/jaddqiu/opsagent/plugins/processors/printer"
	_ "github.com/jaddqiu/opsagent/plugins/processors/regex"
	_ "github.com/jaddqiu/opsagent/plugins/processors/rename"
	_ "github.com/jaddqiu/opsagent/plugins/processors/strings"
	_ "github.com/jaddqiu/opsagent/plugins/processors/topk"
)
