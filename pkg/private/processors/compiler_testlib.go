package processors

import (
	"time"

	"github.com/sshwy/yaoj-core/pkg/private/judger"
	"github.com/sshwy/yaoj-core/pkg/processor"
	"github.com/sshwy/yaoj-core/pkg/utils"
)

// Compile codeforces testlib (or similar) source file using g++.
// For input files, "source" represents source file.
type CompilerTestlib struct {
	// input: source, testlib
	// output: result, log, judgerlog
}

func (r CompilerTestlib) Label() (inputlabel []string, outputlabel []string) {
	return []string{"source", "testlib"}, []string{"result", "log", "judgerlog"}
}

func (r CompilerTestlib) Run(input []string, output []string) *Result {
	if _, err := utils.CopyFile(input[1], "testlib.h"); err != nil {
		return &Result{
			Code: processor.RuntimeError,
			Msg:  "copy: " + err.Error(),
		}
	}
	src := utils.RandomString(10) + ".cpp"
	if _, err := utils.CopyFile(input[0], src); err != nil {
		return &Result{
			Code: processor.RuntimeError,
			Msg:  "copy: " + err.Error(),
		}
	}
	res, err := judger.Judge(
		judger.WithArgument("/dev/null", "/dev/null", output[1], "/usr/bin/g++", src, "-o", output[0]),
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

var _ Processor = CompilerTestlib{}
