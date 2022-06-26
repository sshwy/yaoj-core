package processors

import (
	"fmt"
	"os"
	"strings"

	"github.com/sshwy/yaoj-core/pkg/private/judger"
	"github.com/sshwy/yaoj-core/pkg/processor"
	"github.com/sshwy/yaoj-core/pkg/utils"
)

// Run a program reading from file and print to file and stderr.
// File "config" contains two lines, the first of which acts the same as
// "limit" of RunnerStdio while the second contains two strings denoting input
// file and output file.
type RunnerFileio struct {
	// input: executable, fin, config
	// output: fout, stderr, judgerlog
}

func (r RunnerFileio) Label() (inputlabel []string, outputlabel []string) {
	return []string{"executable", "fin", "config"}, []string{"fout", "stderr", "judgerlog"}
}

func (r RunnerFileio) Run(input []string, output []string) *Result {
	lim, err := os.ReadFile(input[2])
	if err != nil {
		return &Result{
			Code: processor.RuntimeError,
			Msg:  "open config: " + err.Error(),
		}
	}
	lines := strings.Split(string(lim), "\n")
	if len(lines) != 2 {
		return &Result{
			Code: processor.RuntimeError,
			Msg:  "invalid config",
		}
	}
	var inf, ouf string
	fmt.Sscanf(lines[1], "%s%s", &inf, &ouf)
	logger.Printf("inf=%q, out=%q", inf, ouf)
	if _, err := utils.CopyFile(input[1], inf); err != nil {
		return &Result{
			Code: processor.RuntimeError,
			Msg:  "copy: " + err.Error(),
		}
	}
	options := []judger.OptionProvider{
		judger.WithArgument("/dev/null", "/dev/null", output[1], input[0]),
		judger.WithJudger(judger.General),
		judger.WithPolicy("builtin:free"),
		judger.WithLog(output[2], 0, false),
	}
	more, err := parseJudgerLimit(lines[0])
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
	utils.CopyFile(ouf, output[0])
	return res.ProcResult()
}

var _ Processor = RunnerFileio{}
