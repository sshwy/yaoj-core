package processor_test

import (
	"path"
	"testing"

	"github.com/bitfield/script"
	"github.com/sshwy/yaoj-core/judger"
	"github.com/sshwy/yaoj-core/processor"
)

// run `go build -buildmode=plugin` under `example/diff-go` before running this test!
func TestLoad(t *testing.T) {
	proc, err := processor.LoadPlugin("testdata/diff-go/diff-go.so")
	if err != nil {
		t.Error(err)
	}

	t.Log(proc.Label())
}

func TestProcessor(t *testing.T) {
	dir := t.TempDir()
	t.Run("CheckerHcmp", func(t *testing.T) {
		fa := path.Join(dir, "a.hcmp.in")
		fb := path.Join(dir, "b.hcmp.in")
		fc := path.Join(dir, "c.hcmp.out")
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
	})

	t.Run("Compiler", func(t *testing.T) {
		compiler := processor.Compiler{}
		res, err := compiler.Run(
			[]string{"testdata/main.cpp", "testdata/script.sh"},
			[]string{path.Join(dir, "dest"), path.Join(dir, "cp.log"), path.Join(dir, "cpl.judger.log")},
		)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(res)
	})

	t.Run("RunnerStdio", func(t *testing.T) {
		fa := path.Join(dir, "a.runnerstdio.in")
		fb := path.Join(dir, "lim.runnerstdio.in")
		script.Echo("1 2").WriteFile(fa)
		script.Echo("1000 1000 104857600 104857600 104857600 104857600 10").WriteFile(fb)
		runner := processor.RunnerStdio{}
		res, err := runner.Run(
			[]string{path.Join(dir, "dest"), fa, fb},
			[]string{path.Join(dir, "dest.out"), path.Join(dir, "dest.err"), path.Join(dir, "dest.judger.log")},
		)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(res)
		if res.Code != judger.Ok {
			t.Errorf("invalid result")
			return
		}

		output, _ := script.File(path.Join(dir, "dest.out")).String()
		t.Log("output:", output)
	})
}
