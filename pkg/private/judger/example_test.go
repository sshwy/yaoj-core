package judger_test

import (
	"log"
	"path"
	"time"

	"github.com/sshwy/yaoj-core/pkg/private/judger"
)

func ExampleJudge() {
	dir := "/tmp"
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
		log.Fatal(err)
	}
	log.Print(*res)
}
