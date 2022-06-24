package processors

import (
	"os"
	"strings"
	"time"

	"github.com/sshwy/yaoj-core/pkg/private/judger"
	"github.com/sshwy/yaoj-core/pkg/processor"
)

// Execute testlib generator.
// Arguments in "arguments" are seperated by space.
type GeneratorTestlib struct {
	// input: generator arguments
	// output: output stderr judgerlog
}

func (r GeneratorTestlib) Label() (inputlabel []string, outputlabel []string) {
	return []string{"generator", "arguments"}, []string{"output", "stderr", "judgerlog"}
}
func (r GeneratorTestlib) Run(input []string, output []string) *Result {
	args, err := os.ReadFile(input[1])
	if err != nil {
		return &Result{
			Code: processor.RuntimeError,
			Msg:  "open arguments: " + err.Error(),
		}
	}
	argv := strings.Split(string(args), " ")
	finalArgv := []string{"/dev/null", output[0], output[1], input[0]}
	for _, v := range argv {
		if v != "" {
			finalArgv = append(finalArgv, v)
		}
	}
	res, err := judger.Judge(
		judger.WithArgument(finalArgv...),
		judger.WithJudger(judger.General),
		judger.WithPolicy("builtin:free"),
		judger.WithLog(output[2], 0, false),
		judger.WithRealTime(time.Minute),
		judger.WithOutput(10*judger.MB),
	)
	if err != nil {
		return &Result{
			Code: processor.SystemError,
			Msg:  err.Error(),
		}
	}
	return res.ProcResult()
}

var _ Processor = GeneratorTestlib{}
