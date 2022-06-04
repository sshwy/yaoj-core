package processor

import (
	"github.com/sshwy/yaoj-core/judger"
	"github.com/sshwy/yaoj-core/utils"
)

type HashValue = utils.HashValue

type Processor interface {
	Label() (inputlabel []string, outputlabel []string)
	// Given a fixed number of input files, generate output to  corresponding files
	// with execution result. It's ok if result == nil, which means success.
	Run(input []string, output []string) (result *judger.Result, err error)
}
