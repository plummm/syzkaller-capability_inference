package courier

import (
	"sync"

	"github.com/google/syzkaller/pkg/rpctype"
)

const (
	Mutating = 0
	Commands = 1
)

type S2EArgs struct {
	Prog    []byte
	Pointer []byte
}

var MutateArgsQueue = make([]rpctype.ProgQueue, 0)

var CommandsQueue = make([]string, 0)
var Mutex = &sync.Mutex{}

//Append testcase to a queue waits for mutating
func AppendMutatingQueue(p []byte) {
	a := rpctype.ProgQueue{
		Prog: p,
	}
	MutateArgsQueue = append(MutateArgsQueue, a)
}

func AppendCommandsQueue(p []byte) {
	CommandsQueue = append(CommandsQueue, string(p))
}

func RetrieveFirstArg(flag int) interface{} {
	switch flag {
	case Mutating:
		if len(MutateArgsQueue) == 0 {
			break
		}
		p := MutateArgsQueue[0]
		MutateArgsQueue = MutateArgsQueue[1:]
		return p
	case Commands:
		if len(CommandsQueue) == 0 {
			break
		}
		p := CommandsQueue[0]
		CommandsQueue = CommandsQueue[1:]
		return []byte(p)
	}
	return nil
}
