package processors_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/bitfield/script"
	"github.com/sshwy/yaoj-core/pkg/private/processors"
	"github.com/sshwy/yaoj-core/pkg/processor"
)

//go:generate go build -buildmode=plugin -o ./testdata/diff-go ./testdata/diff-go/main.go
func TestLoad(t *testing.T) {
	proc, err := processor.LoadPlugin("testdata/diff-go/main.so")
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
		checker := processors.CheckerHcmp{}
		res := checker.Run([]string{fa, fb}, []string{fc})
		if res.Code != processor.Ok {
			t.Errorf("expect %v, found %v", processor.Ok, res.Code)
			return
		}
		script.Echo("12346").WriteFile(fb)
		res = checker.Run([]string{fa, fb}, []string{fc})
		if res.Code != processor.ExitError {
			t.Errorf("expect %v, found %v", processor.ExitError, res.Code)
			return
		}
	})

	t.Run("Compiler", func(t *testing.T) {
		compiler := processors.Compiler{}
		res := compiler.Run(
			[]string{"testdata/main.cpp", "testdata/script.sh"},
			[]string{path.Join(dir, "dest"), path.Join(dir, "cp.log"), path.Join(dir, "cpl.judger.log")},
		)
		if res.Code != processor.Ok {
			t.Errorf("expect %v, found %v", processor.Ok, res.Code)
			return
		}
		t.Log(res)
	})

	t.Run("RunnerStdio", func(t *testing.T) {
		fa := path.Join(dir, "a.rsi.in")
		fb := path.Join(dir, "lim.rsi.in")
		script.Echo("1 2").WriteFile(fa)
		script.Echo("1000 1000 204857600 204857600 204857600 204857600 10").WriteFile(fb)
		runner := processors.RunnerStdio{}
		res := runner.Run(
			[]string{path.Join(dir, "dest"), fa, fb},
			[]string{path.Join(dir, "dest.out"), path.Join(dir, "dest.err"), path.Join(dir, "dest.judger.log")},
		)
		t.Log(res)
		if res.Code != processor.Ok {
			t.Errorf("invalid result")
			return
		}

		output, _ := script.File(path.Join(dir, "dest.out")).String()
		t.Log("output:", output)
	})

	t.Run("RunnerFileio", func(t *testing.T) {
		script.Exec(fmt.Sprintf("clang++ testdata/main2.cpp -o %s", path.Join(dir, "dest2"))).Wait()
		runner := processors.RunnerFileio{}
		script.Echo("1000 1000 204857600 204857600 204857600 204857600 10\n/tmp/a.in /tmp/a.out").WriteFile(path.Join(dir, "lim2.in"))
		res := runner.Run(
			[]string{path.Join(dir, "dest2"), path.Join(dir, "a.rsi.in"), path.Join(dir, "lim2.in")},
			[]string{path.Join(dir, "dest2.out"), path.Join(dir, "dest2.err"), path.Join(dir, "dest.judger2.log")},
		)
		t.Log(res)
		if res.Code != processor.Ok {
			t.Errorf("invalid result")
			return
		}

		output, _ := script.File(path.Join(dir, "dest2.out")).String()
		t.Log("output:", output)
	})

	t.Run("CheckerTestlib", func(t *testing.T) {
		script.Exec(fmt.Sprintf("clang++ testdata/yesno.cpp -o %s", path.Join(dir, "yesno"))).Wait()
		script.Echo("yes").WriteFile(path.Join(dir, "inp"))
		script.Echo("yes").WriteFile(path.Join(dir, "oup"))
		script.Echo("yes").WriteFile(path.Join(dir, "ans"))
		runner := processors.CheckerTestlib{}
		info, _ := os.Stat(path.Join(dir, "yesno"))
		t.Log(info.Mode())
		res := runner.Run(
			[]string{path.Join(dir, "yesno"), path.Join(dir, "inp"), path.Join(dir, "oup"), path.Join(dir, "ans")},
			[]string{path.Join(dir, "rep"), path.Join(dir, "err.testlib"), path.Join(dir, "jlog.testlib")},
		)
		if res.Code != processor.Ok {
			t.Errorf("invalid result")
			return
		}
		t.Log(res)
		t.Log(script.File(path.Join(dir, "rep")).String())
	})

	t.Run("GeneratorTestlib", func(t *testing.T) {
		script.Exec(fmt.Sprintf("clang++ testdata/igen.cpp -o %s", path.Join(dir, "igen"))).Wait()
		script.Echo("1 4 2 8 5    7").WriteFile(path.Join(dir, "igenparam"))
		runner := processors.GeneratorTestlib{}
		res := runner.Run(
			[]string{path.Join(dir, "igen"), path.Join(dir, "igenparam")},
			[]string{path.Join(dir, "igen.out"), path.Join(dir, "igen.err"), path.Join(dir, "igen.log")},
		)
		if res.Code != processor.Ok {
			t.Errorf("invalid result")
			return
		}
		t.Log(res)
		t.Log(script.File(path.Join(dir, "igen.out")).String())
	})

	t.Run("Inputmaker", func(t *testing.T) {
		script.Echo("raw").WriteFile(path.Join(dir, "rawopt"))
		script.Echo("generator").WriteFile(path.Join(dir, "genopt"))
		runner := processors.Inputmaker{}
		res := runner.Run(
			[]string{path.Join(dir, "igenparam"), path.Join(dir, "rawopt"), "/dev/null"},
			[]string{path.Join(dir, "igen2.out"), path.Join(dir, "igen2.err"), path.Join(dir, "igen2.log")},
		)
		t.Log(res)
		if res.Code != processor.Ok {
			t.Errorf("invalid result")
			return
		}
		t.Log(script.File(path.Join(dir, "igen2.out")).String())
		res = runner.Run(
			[]string{path.Join(dir, "igenparam"), path.Join(dir, "genopt"), path.Join(dir, "igen")},
			[]string{path.Join(dir, "igen3.out"), path.Join(dir, "igen3.err"), path.Join(dir, "igen3.log")},
		)
		t.Log(res)
		if res.Code != processor.Ok {
			t.Errorf("invalid result")
			return
		}
		t.Log(script.File(path.Join(dir, "igen3.out")).String())
	})
	t.Run("manager", func(t *testing.T) {
		mp := processors.GetAll()
		var s = []struct {
			Name   string   `json:"name"`
			Input  []string `json:"input"`
			Output []string `json:"output"`
		}{}
		for k, v := range mp {
			input, output := v.Label()
			s = append(s, struct {
				Name   string   "json:\"name\""
				Input  []string "json:\"input\""
				Output []string "json:\"output\""
			}{
				Name:   k,
				Input:  input,
				Output: output,
			})
		}
		buf, err := json.Marshal(s)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(string(buf))
	})
}
