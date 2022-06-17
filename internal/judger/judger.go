package judger

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var logger = log.New(os.Stderr, "[judger] ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix)

type Option struct {
	Logfile   string
	LogLevel  int
	LogColor  bool
	Policy    string
	PolicyDir string
	Argument  []string
	Environ   []string
	Limit     L
	Runner    Runner
}

type OptionProvider func(*Option)

type ByteValue int64

const KB ByteValue = 1024
const MB ByteValue = KB * KB
const GB ByteValue = KB * MB

func (r ByteValue) String() string {
	num := float64(r)
	if num < 1000 {
		return fmt.Sprint(int64(num), "B")
	} else if num < 1e6 {
		return fmt.Sprintf("%.1f%s", num/1e3, "KB")
	} else if num < 1e9 {
		return fmt.Sprintf("%.1f%s", num/1e6, "MB")
	} else {
		return fmt.Sprintf("%.1f%s", num/1e9, "GB")
	}
}

// Code is required, others are optional
type Result struct {
	// Result status：OK/RE/MLE/...
	Code              StatusCode
	RealTime, CpuTime *time.Duration
	Memory            *ByteValue
	Signal            *int
	Msg               string
}

func (r Result) String() string {
	return fmt.Sprintf("%d{Code: %d, Signal: %d, RealTime: %v, CpuTime: %v, Memory: %v, ErrorMsg: \"%s\"}",
		r.Code, r.Code, *r.Signal, r.RealTime, r.CpuTime, r.Memory, r.Msg)
}

var judgeSync sync.Mutex

func Judge(options ...OptionProvider) (*Result, error) {
	judgeSync.Lock()
	defer judgeSync.Unlock()

	var option = Option{
		Environ:   os.Environ(),
		Policy:    "builtin:free",
		PolicyDir: ".",
		Runner:    General,
		Limit:     make(L),
		Logfile:   "runtime.log",
		LogLevel:  0,
		LogColor:  false,
	}

	for _, v := range options {
		v(&option)
	}

	if err := logSet(option.Logfile, option.LogLevel, option.LogColor); err != nil {
		return nil, err
	}
	defer logClose()

	ctxt := newContext()
	defer ctxt.Free()

	if err := ctxt.SetPolicy(option.PolicyDir, option.Policy); err != nil {
		return nil, err
	}

	if err := ctxt.SetLimit(option.Limit); err != nil {
		return nil, err
	}

	if err := ctxt.SetRunner(option.Argument, option.Environ); err != nil {
		return nil, err
	}

	if err := ctxt.Run(option.Runner); err != nil {
		return nil, err
	}

	result := ctxt.Result()

	return &result, nil
}

// Runners differ in arguments.
//
// For the General: [input] [output] [outerr] [exec] [arguments...]
//
// For the Interactive: [exec] [interactor] [input_itct] [output_itct]
// [outerr_itct] [outerr]. Note that stdin and stdout of interactor and
// executable will be piped together in a two way communication.
func WithArgument(argv ...string) OptionProvider {
	return func(o *Option) {
		o.Argument = argv
	}
}

// default: os.Environ()
func WithEnviron(environ ...string) OptionProvider {
	return func(o *Option) {
		o.Environ = environ
	}
}

// Specify the runner to be used. General (default) or Interactive.
func WithJudger(r Runner) OptionProvider {
	return func(o *Option) {
		o.Runner = r
	}
}

// specify (builtin) policy.
// default: builtin:free
func WithPolicy(name string) OptionProvider {
	return func(o *Option) {
		o.Policy = name
	}
}

func WithPolicyDir(dir string) OptionProvider {
	return func(o *Option) {
		o.PolicyDir = dir
	}
}

// Set real time limitation
func WithRealTime(duration time.Duration) OptionProvider {
	return func(o *Option) {
		o.Limit[realTime] = duration.Milliseconds()
	}
}

func WithCpuTime(duration time.Duration) OptionProvider {
	return func(o *Option) {
		o.Limit[cpuTime] = duration.Milliseconds()
	}
}

// Set virtual memory limitation
func WithVirMemory(space ByteValue) OptionProvider {
	return func(o *Option) {
		o.Limit[virtMem] = int64(space)
	}
}

func WithRealMemory(space ByteValue) OptionProvider {
	return func(o *Option) {
		o.Limit[realMem] = int64(space)
	}
}

func WithStack(space ByteValue) OptionProvider {
	return func(o *Option) {
		o.Limit[stackMem] = int64(space)
	}
}

func WithOutput(space ByteValue) OptionProvider {
	return func(o *Option) {
		o.Limit[outputSize] = int64(space)
	}
}

// Set limitation on number of fileno
func WithFileno(num int) OptionProvider {
	return func(o *Option) {
		o.Limit[filenoLim] = int64(num)
	}
}

// Set logging file. Default is "runtime.log".
func WithLog(file string, level int, color bool) OptionProvider {
	return func(o *Option) {
		o.Logfile = file
		o.LogLevel = level
		o.LogColor = color
	}
}
