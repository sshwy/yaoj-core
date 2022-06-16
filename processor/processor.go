// Package processor provides builtin processors and processor plugin loader
package processor

import (
	"github.com/sshwy/yaoj-core/internal/judger"
)

// Processor takes a series of input (files) and generates a series of outputs.
type Processor interface {
	// Report human-readable label for each input and output
	Label() (inputlabel []string, outputlabel []string)
	// Given a fixed number of input files, generate output to  corresponding files
	// with execution result. It's ok if result == nil, which means success.
	Run(input []string, output []string) (result *judger.Result, err error)
}
