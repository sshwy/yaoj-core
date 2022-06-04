package plugin

import (
	"fmt"
	goPlugin "plugin"

	"github.com/sshwy/yaoj-core/judger"
	"github.com/sshwy/yaoj-core/processor"
)

// Load a go plugin as processor.
// The plugin requires two exported functions:
//
//     func Label() (input []string, output []string)
//     func Main(inputs []string, outputs []string) int
//
func Load(plugin string) (processor.Processor, error) {
	p, err := goPlugin.Open(plugin)
	if err != nil {
		return nil, err
	}

	label, err := p.Lookup("Label")
	if err != nil {
		return nil, err
	}

	var inputLabel, outputLabel []string
	var runner func([]string, []string) int
	if f, ok := label.(func() ([]string, []string)); ok {
		inputLabel, outputLabel = f()
	} else {
		return nil, fmt.Errorf("invalid InputLabel type")
	}

	main, err := p.Lookup("Main")
	if err != nil {
		return nil, err
	}
	if f, ok := main.(func([]string, []string) int); ok {
		runner = f
	} else {
		return nil, fmt.Errorf("invalid Main type")
	}

	return &pluginProcessor{
		inputLabel:  inputLabel,
		outputLabel: outputLabel,
		runner:      runner,
	}, nil
}

type pluginProcessor struct {
	inputLabel, outputLabel []string
	runner                  func([]string, []string) int
}

var _ processor.Processor = (*pluginProcessor)(nil)

func (r *pluginProcessor) Run(input []string, output []string) (result *judger.Result, err error) {
	code := r.runner(input, output)
	if code != 0 {
		return &judger.Result{
			Code:     judger.ExitError,
			ExitCode: code,
		}, nil
	} else {
		return nil, nil
	}
}

func (r *pluginProcessor) Label() (inputlabel []string, outputlabel []string) {
	return r.inputLabel, r.outputLabel
}