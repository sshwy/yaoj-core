package processors

import (
	"os"

	"github.com/sshwy/yaoj-core/pkg/internal/judger"
	"github.com/sshwy/yaoj-core/pkg/processor"
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
func (r RunnerStdio) Run(input []string, output []string) *Result {
	lim, err := os.ReadFile(input[2])
	if err != nil {
		return &Result{
			Code: processor.RuntimeError,
			Msg:  "open limit: " + err.Error(),
		}
	}
	options := []judger.OptionProvider{
		judger.WithArgument(input[1], output[0], output[1], input[0]),
		judger.WithJudger(judger.General),
		judger.WithPolicy("builtin:free"),
		judger.WithLog(output[2], 0, false),
	}
	more, err := parseJudgerLimit(string(lim))
	if err != nil {
		return &Result{
			Code: processor.RuntimeError,
			Msg:  "parse judger limit: " + err.Error(),
		}
	}
	options = append(options, more...)
	res, err := judger.Judge(options...)
	if err != nil {
		return &Result{
			Code: processor.SystemError,
			Msg:  err.Error(),
		}
	}
	return res.ProcResult()
}

var _ Processor = RunnerStdio{}
