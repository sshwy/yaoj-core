package judger_test

import (
	"path"
	"testing"
	"time"

	"github.com/sshwy/yaoj-core/pkg/internal/judger"
)

func TestJudge(t *testing.T) {
	dir := t.TempDir()
	res, err := judger.Judge(
		judger.WithArgument("/dev/null", path.Join(dir, "output"), "/dev/null", "/usr/bin/ls", "."),
		judger.WithJudger(judger.General),
		judger.WithPolicy("builtin:free"),
		judger.WithLog(path.Join(dir, "runtime.log"), 0, false),
		judger.WithRealMemory(3*judger.MB),
		judger.WithStack(3*judger.MB),
		judger.WithVirMemory(3*judger.MB),
		judger.WithRealTime(time.Millisecond*1000),
		judger.WithCpuTime(time.Millisecond*1000),
		judger.WithOutput(3*judger.MB),
		judger.WithFileno(10),
	)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(*res)
}
