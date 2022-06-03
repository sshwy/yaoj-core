package judger

//go:generate go version
//go:generate make -C yaoj-judger

//#cgo CFLAGS: -I./yaoj-judger/include
//#cgo LDFLAGS: -L./yaoj-judger -lyjudger
//#include "./yaoj-judger/include/judger.h"
//#include <stdlib.h>
import "C"
import (
	"fmt"
	"os"
	"sync"
	"time"
)

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

type LimitType int

const (
	RealTime LimitType = C.REAL_TIME
	CpuTime  LimitType = C.CPU_TIME
	// virtual memory
	VirtMem  LimitType = C.VIRTUAL_MEMORY
	RealMem  LimitType = C.ACTUAL_MEMORY
	StackMem LimitType = C.STACK_MEMORY
	// output size
	Output LimitType = C.OUTPUT_SIZE
	Fileno LimitType = C.FILENO
)

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
	coloredres := map[StatusCode]string{
		OK:  "\033[32mOK\033[0m",
		RE:  "\033[31mRE\033[0m",
		MLE: "\033[31mMLE\033[0m",
		TLE: "\033[31mTLE\033[0m",
		OLE: "\033[31mOLE\033[0m",
		SE:  "\033[31mSE\033[0m",
		DSC: "\033[31mDSC\033[0m",
		ECE: "\033[31mECE\033[0m",
	}
	return fmt.Sprintf("%s{Code: %d, Signal: %d, ExitCode: %d, RealTime: %v, CpuTime: %v, Memory: %v}",
		coloredres[r.Code],
		r.Code, r.Signal, r.ExitCode, r.RealTime, r.CpuTime, r.Memory)
}

type StatusCode int

const (
	OK  StatusCode = C.OK
	RE  StatusCode = C.RE
	MLE StatusCode = C.MLE
	TLE StatusCode = C.TLE
	OLE StatusCode = C.OLE
	SE  StatusCode = C.SE
	DSC StatusCode = C.DSC
	ECE StatusCode = C.ECE
)

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
		o.Limit[RealTime] = duration.Milliseconds()
	}
}

func WithCpuTime(duration time.Duration) OptionProvider {
	return func(o *Option) {
		o.Limit[CpuTime] = duration.Milliseconds()
	}
}

func WithVirMemory(space ByteValue) OptionProvider {
	return func(o *Option) {
		o.Limit[VirtMem] = int64(space)
	}
}

func WithRealMemory(space ByteValue) OptionProvider {
	return func(o *Option) {
		o.Limit[RealMem] = int64(space)
	}
}

func WithStack(space ByteValue) OptionProvider {
	return func(o *Option) {
		o.Limit[StackMem] = int64(space)
	}
}

func WithOutput(space ByteValue) OptionProvider {
	return func(o *Option) {
		o.Limit[Output] = int64(space)
	}
}

func WithFileno(num int) OptionProvider {
	return func(o *Option) {
		o.Limit[Fileno] = int64(num)
	}
}
