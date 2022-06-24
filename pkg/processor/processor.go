// Package processor provides builtin processors and processor plugin loader
package processor

import (
	"time"

	"github.com/sshwy/yaoj-core/pkg/utils"
)

// Processor takes a series of input (files) and generates a series of outputs.
type Processor interface {
	// Report human-readable label for each input and output
	Label() (inputlabel []string, outputlabel []string)
	// Given a fixed number of input files, generate output to  corresponding files
	// with execution result. It's ok if result == nil, which means success.
	Run(input []string, output []string) (result *Result)
}

type Code int

const (
	Ok Code = iota
	RuntimeError
	MemoryExceed
	TimeExceed
	OoutputExceed
	SystemError
	DangerousSyscall
	ExitError
)

// Code is required, others are optional
type Result struct {
	// Result status：OK/RE/MLE/...
	Code              Code
	RealTime, CpuTime *time.Duration
	Memory            *utils.ByteValue
	Msg               string
}
