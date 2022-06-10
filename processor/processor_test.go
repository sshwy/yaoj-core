package processor_test

import (
	"path"
	"testing"

	"github.com/bitfield/script"
	"github.com/sshwy/yaoj-core/judger"
	"github.com/sshwy/yaoj-core/processor"
)

func TestCheckerHcmp(t *testing.T) {
	dir := t.TempDir()
	fa := path.Join(dir, "a.in")
	fb := path.Join(dir, "b.in")
	fc := path.Join(dir, "c.out")
	script.Echo("12345").WriteFile(fa)
	script.Echo("12345").WriteFile(fb)
	checker := processor.CheckerHcmp{}
	res, err := checker.Run([]string{fa, fb}, []string{fc})
	if err != nil {
		t.Error(err)
		return
	}
	if res.Code != judger.Ok {
		t.Errorf("expect %v, found %v", judger.Ok, res.Code)
		return
	}
	script.Echo("12346").WriteFile(fb)
	res, err = checker.Run([]string{fa, fb}, []string{fc})
	if err != nil {
		t.Error(err)
		return
	}
	if res.Code != judger.ExitError {
		t.Errorf("expect %v, found %v", judger.ExitError, res.Code)
		return
	}
}

// run `go build -buildmode=plugin` under `example/diff-go` before running this test!
func TestLoad(t *testing.T) {
	proc, err := processor.LoadPlugin("example/diff-go/diff-go.so")
	if err != nil {
		t.Error(err)
	}

	t.Log(proc.Label())
}
