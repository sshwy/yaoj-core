package workflow_test

import (
	"testing"

	"github.com/k0kubun/pp/v3"
	"github.com/sshwy/yaoj-core/pkg/workflow"
)

func TestWorkflow(t *testing.T) {
	var b workflow.Builder
	b.SetNode("compile", "compiler", false)
	b.SetNode("run", "runner:stdio", true)
	b.SetNode("check", "checker:hcmp", false)
	b.AddInbound("submission", "source", "compile", "source")
	b.AddInbound("static", "compilescript", "compile", "script")
	b.AddInbound("static", "limitation", "run", "limit")
	b.AddInbound("testcase", "input", "run", "stdin")
	b.AddInbound("testcase", "answer", "check", "ans")
	b.AddEdge("compile", "result", "run", "executable")
	b.AddEdge("run", "stdout", "check", "out")
	graph, err := b.WorkflowGraph()
	if err != nil {
		t.Error(err)
		return
	}
	pp.Print(graph)

	_ = workflow.Workflow{
		WorkflowGraph: graph,
		Analyzer:      workflow.DefaultAnalyzer{},
	}
	//dir := t.TempDir()
	//res, err := workflow.Run(w, dir, map[string]*map[string]string{
	//	"testcase": {
	//		"input":  "testdata/main.in",
	//		"answer": "testdata/main.ans",
	//	},
	//	"static": {
	//		"limitation":    "testdata/main.lim",
	//		"compilescript": "testdata/script.sh",
	//	},
	//	"submission": {
	//		"source": "testdata/main.cpp",
	//	},
	//}, 100)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//if res.Score != res.Fullscore {
	//	t.Errorf("score=%f, expect %f", res.Score, res.Fullscore)
	//	return
	//}
	//t.Log(pp.Sprint(*res))
	//t.Log(string(w.Serialize()))
	//w2, err := workflow.Load(w.Serialize())
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//t.Log(string(w2.Serialize()))
}

//go:generate go build -buildmode=plugin -o ./testdata ./testdata/custom_analyzer.go
func TestLoadAnalyzer(t *testing.T) {
	a, err := workflow.LoadAnalyzer("testdata/custom_analyzer.so")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(pp.Sprint(a))
}
