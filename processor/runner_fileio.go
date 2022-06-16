package processor

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sshwy/yaoj-core/internal/judger"
)

// Run a program reading from file and print to file and stderr.
// File "config" contains two lines, the first of which acts the same as
// "limit" of RunnerStdio while the second contains two strings denoting input
// file and output file.
type RunnerFileio struct {
	// input: executable, fin, config
	// output: fout, stderr, judgerlog
}

func (r RunnerFileio) Label() (inputlabel []string, outputlabel []string) {
	return []string{"executable", "fin", "config"}, []string{"fout", "stderr", "judgerlog"}
}

func (r RunnerFileio) Run(input []string, output []string) (result *judger.Result, err error) {
	lim, err := os.ReadFile(input[2])
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(lim), "\n")
	if len(lines) != 2 {
		return nil, fmt.Errorf("invalid config")
	}
	var inf, ouf string
	fmt.Sscanf(lines[1], "%s%s", &inf, &ouf)
	if _, err := copyFile(input[1], inf); err != nil {
		return nil, err
	}
	options := []judger.OptionProvider{
		judger.WithArgument("/dev/null", "/dev/null", output[1], input[0]),
		judger.WithJudger(judger.General),
		judger.WithPolicy("builtin:free"),
		judger.WithLog(output[2], 0, false),
	}
	more, err := parseJudgerLimit(lines[0])
	if err != nil {
		return nil, err
	}
	options = append(options, more...)
	res, err := judger.Judge(options...)
	if err != nil {
		return nil, err
	}
	copyFile(ouf, output[0])
	return res, nil
}

var _ Processor = RunnerFileio{}

func copyFile(src, dst string) (int64, error) {
	if src == dst {
		return 0, fmt.Errorf("same path")
	}
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
