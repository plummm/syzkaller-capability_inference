package courier

import (
	"github.com/google/syzkaller/prog"
)

var AnalyzerPath string
var ConfirmedSuccess = false

func AppendTestcase(p *prog.Prog) {
	AppendMutatingQueue(p.Serialize())
}
