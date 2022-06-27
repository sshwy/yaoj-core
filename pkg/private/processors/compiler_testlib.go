package processors

import (
	"os"
	"time"

	_ "embed"

	"github.com/sshwy/yaoj-core/pkg/private/judger"
	"github.com/sshwy/yaoj-core/pkg/processor"
	"github.com/sshwy/yaoj-core/pkg/utils"
)

//go:embed testlib.h
var testlib []byte

// Compile codeforces testlib source file using g++.
// For input files, "source" represents source file.
type CompilerTestlib struct {
	// input: source
	// output: result, log, judgerlog
}

func (r CompilerTestlib) Label() (inputlabel []string, outputlabel []string) {
	return []string{"source"}, []string{"result", "log", "judgerlog"}
}

func (r CompilerTestlib) Run(input []string, output []string) *Result {
	file, err := os.Create("testlib.h")
	if err != nil {
		return &Result{
			Code: processor.RuntimeError,
			Msg:  "open: " + err.Error(),
		}
	}
	_, err = file.Write(testlib)
	if err != nil {
		return &Result{
			Code: processor.RuntimeError,
			Msg:  "write: " + err.Error(),
		}
	}
	file.Close()

	src := utils.RandomString(10) + ".cpp"
	if _, err := utils.CopyFile(input[0], src); err != nil {
		return &Result{
			Code: processor.RuntimeError,
			Msg:  "copy: " + err.Error(),
		}
	}
	res, err := judger.Judge(
		judger.WithArgument("/dev/null", "/dev/null", output[1], "/usr/bin/g++", src, "-o", output[0], "-O2", "-Wall"),
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
