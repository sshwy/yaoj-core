package workflow_test

import (
	"testing"

	"github.com/k0kubun/pp/v3"
	"github.com/sshwy/yaoj-core/internal/judger"
	"github.com/sshwy/yaoj-core/workflow"
)

type testanalyzer struct {
}

func (r testanalyzer) Analyze(nodes []workflow.RuntimeNode, score float64) workflow.Result {
	display := []workflow.ResultFileDisplay{
		workflow.FileDisplay(nodes[1].Output[0], "output", 50),
		workflow.FileDisplay(nodes[1].Output[1], "error", 50),
	}
	if nodes[2].Result.Code == judger.Ok {
		return workflow.Result{
			Fullscore: score,
			Score:     score,
			File:      display,
		}
	}
	return workflow.Result{
		Fullscore: score,
		Score:     0,
		File:      display,
	}
}

func TestWorkflow(t *testing.T) {
	w := workflow.Workflow{
		WorkflowGraph: &workflow.WorkflowGraph{
			Edge: []workflow.Edge{
				{
					From: workflow.Outbound{
						Index:      0, // compiler
						LabelIndex: 0, // result
					},
					To: workflow.Inbound{
						Index:      1, // runner
						LabelIndex: 0, // executable
					},
				},
				{
					From: workflow.Outbound{
						Index:      1, // runner
						LabelIndex: 0, // stdout
					},
					To: workflow.Inbound{
						Index:      2, // checker
						LabelIndex: 0, // out
					},
				},
			},
			Node: []workflow.Node{
				{ProcName: "compiler"},
				{ProcName: "runner:stdio"},
				{ProcName: "checker:hcmp"},
			},
			Inbound: map[string]*map[string][]workflow.Inbound{
				"testcase": {
					"input": {{
						Index:      1, // runner:stdio
						LabelIndex: 1, // stdin
					}},
					"answer": {{
						Index:      2, // checker
						LabelIndex: 1, // ans
					}},
				},
				"option": {
					"limitation": {{
						Index:      1, // runner:stdio
						LabelIndex: 2, // limit
					}},
					"compilescript": {{
						Index:      0, // compiler
						LabelIndex: 1, // script
					}},
				},
				"submission": {
					"source": {{
						Index:      0, // compiler
						LabelIndex: 0, // source
					}},
				},
			},
		},
		Analyzer: testanalyzer{},
	}
	dir := t.TempDir()
	res, err := workflow.Run(w, dir, map[string]*map[string]string{
		"testcase": {
			"input":  "testdata/main.in",
			"answer": "testdata/main.ans",
		},
		"option": {
			"limitation":    "testdata/main.lim",
			"compilescript": "testdata/script.sh",
		},
		"submission": {
			"source": "testdata/main.cpp",
		},
	}, 100)
	if err != nil {
		t.Error(err)
		return
	}
	if res.Score != res.Fullscore {
		t.Errorf("score=%f, expect %f", res.Score, res.Fullscore)
		return
	}
	t.Log(pp.Sprint(*res))
	t.Log(string(w.Serialize()))
	w2, err := workflow.Load(w.Serialize())
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(w2.Serialize()))
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
