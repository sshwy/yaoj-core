package judger_test

import (
	"testing"

	"github.com/sshwy/yaoj-worker/judger"
)

func TestJudge(t *testing.T) {
	res, err := judger.Judge(
		judger.WithArgument("/dev/null", "output", "/dev/null", "/usr/bin/ls", "."),
		judger.WithJudger(judger.General),
		judger.WithPolicy("builtin:free"),
		judger.WithRealMemory(3*judger.MB),
	)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(*res)
}
