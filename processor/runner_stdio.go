package processor

import (
	"github.com/bitfield/script"
	"github.com/sshwy/yaoj-core/judger"
)

// Run a program reading from stdin and print to stdout and stderr.
// For "limit", it contains a series of number seperated by space, denoting
// real time (ms), cpu time (ms), virtual memory (byte), real memory (byte),
// stack memory (byte), output limit (byte), fileno limitation respectively.
type RunnerStdio struct {
	// input: executable, stdin, limit
	// output: stdout, stderr, judgerlog
}

func (r RunnerStdio) Label() (inputlabel []string, outputlabel []string) {
	return []string{"executable", "stdin", "limit"}, []string{"stdout", "stderr", "judgerlog"}
}
func (r RunnerStdio) Run(input []string, output []string) (result *judger.Result, err error) {
	lim, err := script.File(input[2]).String()
	if err != nil {
		return nil, err
	}
	options := []judger.OptionProvider{
		judger.WithArgument(input[1], output[0], output[1], input[0]),
		judger.WithJudger(judger.General),
		judger.WithPolicy("builtin:free"),
		judger.WithLog(output[2], 0, false),
	}
	options = append(options, parseJudgerLimit(lim)...)
	res, err := judger.Judge(options...)
	if err != nil {
		return nil, err
	}
	return res, nil
}

var _ Processor = RunnerStdio{}
