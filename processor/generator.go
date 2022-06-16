package processor

import (
	"os"
	"strings"

	"github.com/sshwy/yaoj-core/internal/judger"
)

// Generate input files from either raw content or testlib generator.
// If "option" contains "_raw" substring, then source should be a text file,
// otherwise a testlib generator is expected.
type Generator struct {
	// source option
	// output: output stderr judgerlog
}

func (r Generator) Label() (inputlabel []string, outputlabel []string) {
	return []string{"source", "option"}, []string{"output", "stderr", "judgerlog"}
}

func (r Generator) Run(input []string, output []string) (result *judger.Result, err error) {
	args, err := os.ReadFile(input[1])
	if err != nil {
		return nil, err
	}
	if strings.Contains(string(args), "_raw") {
		if _, err := copyFile(input[0], output[0]); err != nil {
			return nil, err
		}
		return &judger.Result{Code: judger.Ok}, nil
	} else { // testlib
		runner := GeneratorTestlib{}
		return runner.Run(input, output)
	}
}

var _ Processor = Generator{}
