package judger

import (
	"errors"
	"fmt"
	"time"
	"unsafe"
)

//go:generate go version
//go:generate make -C yaoj-judger

//#cgo CFLAGS: -I./yaoj-judger/include
//#cgo LDFLAGS: -L./yaoj-judger -lyjudger
//#include "./yaoj-judger/include/judger.h"
//#include <stdlib.h>
import "C"

// StatusCode describes final status of a execution. Note that Ok is
// the only status for success
type StatusCode int

const (
	Ok               StatusCode = C.OK
	RuntimeError     StatusCode = C.RE
	MemoryExceed     StatusCode = C.MLE
	TimeExceed       StatusCode = C.TLE
	OoutputExceed    StatusCode = C.OLE
	SystemError      StatusCode = C.SE
	DangerousSyscall StatusCode = C.DSC
	ExitError        StatusCode = C.ECE
)

type LimitType int

const (
	realTime LimitType = C.REAL_TIME
	cpuTime  LimitType = C.CPU_TIME
	// virtual memory
	virtMem  LimitType = C.VIRTUAL_MEMORY
	realMem  LimitType = C.ACTUAL_MEMORY
	stackMem LimitType = C.STACK_MEMORY
	// output size
	outputSize LimitType = C.OUTPUT_SIZE
	filenoLim  LimitType = C.FILENO
)

func boolToInt(v bool) int {
	if v {
		return 1
	} else {
		return 0
	}
}

// Set logging options.
// MUST be executed before creating context.
//
// filename set perform log file.
// log_level determine minimum log level (DEBUG, INFO, WARN, ERROR = 0, 1, 2, 3)
// with_color whether use ASCII color controller character
func logSet(filename string, level int, color bool) error {
	logger.Printf("LogSet: file=%s level=%d color=%v", filename, level, color)
	var cfilename *C.char = C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))

	res := C.log_set(cfilename, C.int(level), C.int(boolToInt(color)))
	if res != 0 {
		return errors.New("log_set return non zero")
	}
	return nil
}

// Close log file.
func logClose() {
	C.log_close()
}

type context struct {
	ctxt C.yjudger_ctxt_t
}

func newContext() context {
	return context{ctxt: C.yjudger_ctxt_create()}
}

func (r context) Result() Result {
	result := C.yjudger_result(r.ctxt)
	signal := int(result.signal)
	exitCode := int(result.exit_code)
	realTime := time.Duration(int(result.real_time) * int(time.Millisecond))
	cpuTime := time.Duration(int(result.cpu_time) * int(time.Millisecond))
	memory := ByteValue(result.real_memory)

	return Result{
		Code:     StatusCode(result.code),
		Signal:   &signal,
		Msg:      fmt.Sprintf("Exit with code %d", exitCode),
		RealTime: &realTime,
		CpuTime:  &cpuTime,
		Memory:   &memory,
	}
}

func (r context) Free() {
	C.yjudger_ctxt_free(r.ctxt)
}

func (r context) SetPolicy(dirname string, policy string) error {
	var cdirname, cpolicy *C.char = C.CString(dirname), C.CString(policy)
	defer C.free(unsafe.Pointer(cdirname))
	defer C.free(unsafe.Pointer(cpolicy))

	flag := C.yjudger_set_policy(r.ctxt, cdirname, cpolicy)
	if flag != 0 {
		return errors.New("set policy error")
	}
	return nil
}

func (r context) SetBuiltinPolicy(policy string) error {
	return r.SetPolicy(".", "builtin:"+policy)
}

func cCharArray(a []string) []*C.char {
	var ca []*C.char = make([]*C.char, len(a)+1)
	for i := range a {
		ca[i] = C.CString(a[i])
	}
	ca[len(ca)-1] = nil
	return ca
}

func cFreeCharArray(ca []*C.char) {
	for _, val := range ca {
		if val != nil {
			C.free(unsafe.Pointer(val))
		}
	}
}

func (r context) SetRunner(argv []string, env []string) error {
	cargv, cenv := cCharArray(argv), cCharArray(env)
	defer cFreeCharArray(cargv)
	defer cFreeCharArray(cenv)

	flag := C.yjudger_set_runner(r.ctxt, C.int(len(argv)), &cargv[0], &cenv[0])
	if flag != 0 {
		return errors.New("set runner error")
	}
	return nil
}

type Runner int

// Runner type
const (
	General     Runner = 0
	Interactive Runner = 1
)

func (r context) Run(runner Runner) error {
	var flag C.int
	switch runner {
	case General:
		flag = C.yjudger_general(r.ctxt)
	case Interactive:
		flag = C.yjudger_interactive(r.ctxt)
	default:
		return errors.New("unknown runner: " + fmt.Sprint(runner))
	}
	if flag != 0 {
		return errors.New("perform general error")
	}
	return nil
}

func (r context) RunFork(runner Runner) Result {
	result := C.yjudger_general_fork(r.ctxt)
	signal := int(result.signal)
	exitCode := int(result.exit_code)
	realTime := time.Duration(int(result.real_time) * int(time.Millisecond))
	cpuTime := time.Duration(int(result.cpu_time) * int(time.Millisecond))
	memory := ByteValue(result.real_memory)

	return Result{
		Code:     StatusCode(result.code),
		Signal:   &signal,
		Msg:      fmt.Sprintf("Exit with code %d", exitCode),
		RealTime: &realTime,
		CpuTime:  &cpuTime,
		Memory:   &memory,
	}
}

// short cut for Limitation
type L map[LimitType]int64

func (r context) SetLimit(options L) error {
	for key, val := range options {
		switch key {
		case realTime:
			C.yjudger_set_limit(r.ctxt, C.REAL_TIME, C.int(val))
		case cpuTime:
			C.yjudger_set_limit(r.ctxt, C.CPU_TIME, C.int(val))
		case virtMem:
			C.yjudger_set_limit(r.ctxt, C.VIRTUAL_MEMORY, C.int(val))
		case realMem:
			C.yjudger_set_limit(r.ctxt, C.ACTUAL_MEMORY, C.int(val))
		case stackMem:
			C.yjudger_set_limit(r.ctxt, C.STACK_MEMORY, C.int(val))
		case outputSize:
			C.yjudger_set_limit(r.ctxt, C.OUTPUT_SIZE, C.int(val))
		case filenoLim:
			C.yjudger_set_limit(r.ctxt, C.FILENO, C.int(val))
		default:
			return fmt.Errorf("unknown limit type: %d", key)
		}
	}
	return nil
}
