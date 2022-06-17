package processor

import (
	"os"
	"time"

	"github.com/sshwy/yaoj-core/internal/judger"
)

// Compile source file in all language.
// For input files, "source" represents source file, "script" represents
// bash script to compile, where $1 gives source file path and $2 gives output file path
type Compiler struct {
	// input: source, script
	// output: result, log, judgerlog
}

func (r Compiler) Label() (inputlabel []string, outputlabel []string) {
	return []string{"source", "script"}, []string{"result", "log", "judgerlog"}
}

func (r Compiler) Run(input []string, output []string) *judger.Result {
	if err := os.Chmod(input[1], 0744); err != nil { // -rwxr--r--
		return &judger.Result{
			Code: judger.RuntimeError,
			Msg:  "open script: " + err.Error(),
		}
	}
	defer os.Chmod(input[1], 0644) // -rw-r--r--
	res, err := judger.Judge(
		judger.WithArgument("/dev/null", "/dev/null", output[1], input[1], input[0], output[0]),
		judger.WithJudger(judger.General),
		judger.WithPolicy("builtin:free"),
		judger.WithLog(output[2], 0, false),
		judger.WithRealTime(time.Minute),
		judger.WithOutput(10*judger.MB),
	)
	if err != nil {
		return &judger.Result{
			Code: judger.SystemError,
			Msg:  err.Error(),
		}
	}
	return res
}

var _ Processor = Compiler{}
