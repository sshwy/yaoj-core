package processor

import (
	"os"
	"strings"

	"github.com/sshwy/yaoj-core/pkg/internal/judger"
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

func (r Inputmaker) Run(input []string, output []string) *judger.Result {
	option, err := os.ReadFile(input[1])
	if err != nil {
		return &judger.Result{
			Code: judger.RuntimeError,
			Msg:  "open option: " + err.Error(),
		}
	}
	if strings.Contains(string(option), "raw") {
		if _, err := utils.CopyFile(input[0], output[0]); err != nil {
			return &judger.Result{
				Code: judger.RuntimeError,
				Msg:  "copy: " + err.Error(),
			}
		}
		return &judger.Result{Code: judger.Ok}
	} else { // testlib
		runner := GeneratorTestlib{}
		return runner.Run([]string{input[2], input[0]}, output)
	}
}

var _ Processor = Inputmaker{}
