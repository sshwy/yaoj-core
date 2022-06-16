package processor

import (
	"fmt"
	"os"

	"github.com/sshwy/yaoj-core/internal/judger"
)

// Compares two signed huge (big) integers.
// Validates that both integers (in the output and in the answer) are well-formatted.
type CheckerHcmp struct {
}

func (r CheckerHcmp) Label() (inputlabel []string, outputlabel []string) {
	return []string{"out", "ans"}, []string{"result"}
}
func (r CheckerHcmp) Run(input []string, output []string) (*judger.Result, error) {
	filea, err := os.Open(input[0])
	if err != nil {
		return nil, err
	}
	defer filea.Close()
	fileb, err := os.Open(input[1])
	if err != nil {
		return nil, err
	}
	defer fileb.Close()

	var a, b string
	fmt.Fscanf(filea, "%s", &a)
	fmt.Fscanf(fileb, "%s", &b)

	if a == b {
		return &judger.Result{
			Code: judger.Ok,
		}, nil
	} else {
		filec, err := os.OpenFile(output[0], os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0744)
		if err != nil {
			return nil, err
		}
		defer filec.Close()

		fmt.Fprintf(filec, "wa: expected '%s', found '%s'", b, a)
		var exitcode = 1

		return &judger.Result{
			Code:     judger.ExitError,
			ExitCode: &exitcode,
		}, nil
	}
}

var _ Processor = CheckerHcmp{}
