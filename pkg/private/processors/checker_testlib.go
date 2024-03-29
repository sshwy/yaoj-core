package processors

import (
	"time"

	"github.com/sshwy/yaoj-core/pkg/private/judger"
	"github.com/sshwy/yaoj-core/pkg/processor"
)

// Execute testlib checker
type CheckerTestlib struct {
	// input: checker input output answer
	// output: xmlreport stderr judgerlog
}

func (r CheckerTestlib) Label() (inputlabel []string, outputlabel []string) {
	return []string{"checker", "input", "output", "answer"},
		[]string{"xmlreport", "stderr", "judgerlog"}
}
func (r CheckerTestlib) Run(input []string, output []string) *Result {
	res, err := judger.Judge(
		judger.WithArgument("/dev/null", "/dev/null", output[1], input[0],
			input[1], input[2], input[3], output[0], "-appes"),
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

var _ Processor = CheckerTestlib{}
