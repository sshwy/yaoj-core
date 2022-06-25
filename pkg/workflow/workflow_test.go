package workflow_test

import (
	"testing"

	"github.com/k0kubun/pp/v3"
	"github.com/sshwy/yaoj-core/pkg/workflow"
)

//go:generate go build -buildmode=plugin -o ./testdata ./testdata/custom_analyzer.go
func TestLoadAnalyzer(t *testing.T) {
	a, err := workflow.LoadAnalyzer("testdata/custom_analyzer.so")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(pp.Sprint(a))
}
