package processor

import (
	"fmt"
	"time"

	"github.com/bitfield/script"
	"github.com/sshwy/yaoj-core/judger"
)

// Run program reading from stdin and print to stdout and stderr.
// For "limit", it contains a series of number seperated by space, denoting
// real time (ms), cpu time (ms), virtual memory (byte), real memory (byte),
// stack memory (byte), output limit (byte), fileno limitation respectively.
type RunnerStdio struct {
}

func (r RunnerStdio) Label() (inputlabel []string, outputlabel []string) {
	return []string{"executable", "stdin", "limit"}, []string{"stdout", "stderr", "judgerlog"}
}
func (r RunnerStdio) Run(input []string, output []string) (result *judger.Result, err error) {
	lim, err := script.File(input[2]).String()
	if err != nil {
		return nil, err
	}
	var rt, ct, vm, rm, sm, ol, fl int
	fmt.Sscanf(lim, "%d%d%d%d%d%d%d", &rt, &ct, &vm, &rm, &sm, &ol, &fl)
	options := []judger.OptionProvider{
		judger.WithArgument(input[1], output[0], output[1], input[0]),
		judger.WithJudger(judger.General),
		judger.WithPolicy("builtin:free"),
		judger.WithLog(output[2], 0, false),
	}
	if rt > 0 {
		options = append(options, judger.WithRealTime(time.Millisecond*time.Duration(rt)))
	}
	if ct > 0 {
		options = append(options, judger.WithCpuTime(time.Millisecond*time.Duration(ct)))
	}
	if vm > 0 {
		options = append(options, judger.WithVirMemory(judger.ByteValue(vm)))
	}
	if rm > 0 {
		options = append(options, judger.WithRealMemory(judger.ByteValue(rm)))
	}
	if sm > 0 {
		options = append(options, judger.WithStack(judger.ByteValue(sm)))
	}
	if ol > 0 {
		options = append(options, judger.WithOutput(judger.ByteValue(ol)))
	}
	if fl > 0 {
		options = append(options, judger.WithFileno(fl))
	}
	res, err := judger.Judge(options...)
	if err != nil {
		return nil, err
	}
	return res, nil
}

var _ Processor = RunnerStdio{}
