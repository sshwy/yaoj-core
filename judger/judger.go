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

type Result struct {
	// 结果状态：OK/RE/MLE/...
	Code              StatusCode
	RealTime, CpuTime time.Duration
	Memory            ByteValue
	Signal, ExitCode  int
}

func (r Result) String() string {
	return fmt.Sprintf("%d{Code: %d, Signal: %d, ExitCode: %d, RealTime: %v, CpuTime: %v, Memory: %v}",
		r.Code, r.Code, r.Signal, r.ExitCode, r.RealTime, r.CpuTime, r.Memory)
}

// 用于同步操作
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

	if err := LogSet(option.Logfile, option.LogLevel, option.LogColor); err != nil {
		return nil, err
	}
	defer LogClose()

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

func WithFileno(num int) OptionProvider {
	return func(o *Option) {
		o.Limit[filenoLim] = int64(num)
	}
}

func WithLog(file string, level int, color bool) OptionProvider {
	return func(o *Option) {
		o.Logfile = file
		o.LogLevel = level
		o.LogColor = color
	}
}
