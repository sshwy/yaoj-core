package processors

import (
	"fmt"
	"os"

	"github.com/sshwy/yaoj-core/pkg/private/judger"
	"github.com/sshwy/yaoj-core/pkg/processor"
)

// Compares two signed huge (big) integers.
// Validates that both integers (in the output and in the answer) are well-formatted.
type CheckerHcmp struct {
	// input: out, ans
	// output: result
}

func (r CheckerHcmp) Label() (inputlabel []string, outputlabel []string) {
	return []string{"out", "ans"}, []string{"result"}
}
func (r CheckerHcmp) Run(input []string, output []string) *processor.Result {
	filea, err := os.Open(input[0])
	if err != nil {
		return &Result{
			Code: processor.RuntimeError,
			Msg:  fmt.Sprintf("open (out) %s: %s", input[0], err.Error()),
		}
	}
	defer filea.Close()
	fileb, err := os.Open(input[1])
	if err != nil {
		return &Result{
			Code: processor.RuntimeError,
			Msg:  fmt.Sprintf("open (ans) %s: %s", input[1], err.Error()),
		}
	}
	defer fileb.Close()

	var a, b string
	fmt.Fscanf(filea, "%s", &a)
	fmt.Fscanf(fileb, "%s", &b)

	if a == b {
		return (&judger.Result{
			Code: judger.Ok,
		}).ProcResult()
	} else {
		filec, err := os.OpenFile(output[0], os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0744)
		if err != nil {
			return &Result{
				Code: processor.RuntimeError,
				Msg:  fmt.Sprintf("open (result) %s: %s", output[0], err.Error()),
			}
		}
		defer filec.Close()

		fmt.Fprintf(filec, "wa: expected '%s', found '%s'", b, a)

		return &Result{
			Code: processor.ExitError,
			Msg:  "exit with code 1",
		}
	}
}

var _ processor.Processor = CheckerHcmp{}
