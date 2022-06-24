package processors

import (
	"os"
	"strings"

	"github.com/sshwy/yaoj-core/pkg/processor"
	"github.com/sshwy/yaoj-core/pkg/utils"
)

// Inputmaker make input according to "option": "raw" means "source" provides
// input content, "generator" means execute "generator" with arguments in
// "source", separated by space.
type Inputmaker struct {
	// source option generator
	// output: result stderr judgerlog
}

func (r Inputmaker) Label() (inputlabel []string, outputlabel []string) {
	return []string{"source", "option", "generator"}, []string{"result", "stderr", "judgerlog"}
}

func (r Inputmaker) Run(input []string, output []string) *Result {
	option, err := os.ReadFile(input[1])
	if err != nil {
		return &Result{
			Code: processor.RuntimeError,
			Msg:  "open option: " + err.Error(),
		}
	}
	if strings.Contains(string(option), "raw") {
		if _, err := utils.CopyFile(input[0], output[0]); err != nil {
			return &Result{
				Code: processor.RuntimeError,
				Msg:  "copy: " + err.Error(),
			}
		}
		return &Result{Code: processor.Ok}
	} else { // testlib
		runner := GeneratorTestlib{}
		return runner.Run([]string{input[2], input[0]}, output)
	}
}

var _ Processor = Inputmaker{}
