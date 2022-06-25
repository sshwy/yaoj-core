package test_test

import (
	"testing"

	"github.com/k0kubun/pp/v3"
	"github.com/sshwy/yaoj-core/pkg/workflow"
)

var wkflGraph workflow.WorkflowGraph

func MakeWorkflowGraph(t *testing.T) {
	var b workflow.Builder
	b.SetNode("compile", "compiler", false)
	b.SetNode("run", "runner:stdio", true)
	b.SetNode("check", "checker:hcmp", false)
	b.AddInbound(workflow.Gsubm, "source", "compile", "source")
	b.AddInbound(workflow.Gstatic, "compilescript", "compile", "script")
	b.AddInbound(workflow.Gstatic, "limitation", "run", "limit")
	b.AddInbound(workflow.Gtests, "input", "run", "stdin")
	b.AddInbound(workflow.Gtests, "answer", "check", "ans")
	b.AddEdge("compile", "result", "run", "executable")
	b.AddEdge("run", "stdout", "check", "out")
	graph, err := b.WorkflowGraph()
	if err != nil {
		t.Error(err)
		return
	}
	pp.Print(graph)
	wkflGraph = *graph

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
