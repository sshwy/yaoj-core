package plugin_test

import (
	"testing"

	"github.com/sshwy/yaoj-core/processor/plugin"
)

// run `go build -buildmode=plugin` under `example/diff-go` before running this test!
func TestLoad(t *testing.T) {
	proc, err := plugin.Load("example/diff-go/diff-go.so")
	if err != nil {
		t.Error(err)
	}

	t.Log(proc.Label())
}
