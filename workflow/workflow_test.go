package workflow_test

import (
	"testing"

	"github.com/k0kubun/pp/v3"
	"github.com/sshwy/yaoj-core/workflow"
)

func TestWorkflow(t *testing.T) {
	w := workflow.Workflow{
		WorkflowGraph: &workflow.WorkflowGraph{
			Edge: []workflow.Edge{
				{
					From: workflow.Outbound{
						Name:       "compile",
						LabelIndex: 0, // result
					},
					To: workflow.Inbound{
						Name:       "run",
						LabelIndex: 0, // executable
					},
				},
				{
					From: workflow.Outbound{
						Name:       "run",
						LabelIndex: 0, // stdout
					},
					To: workflow.Inbound{
						Name:       "check",
						LabelIndex: 0, // out
					},
				},
			},
			Node: map[string]workflow.Node{
				"check":   {ProcName: "checker:hcmp"},
				"run":     {ProcName: "runner:stdio", Key: true},
				"compile": {ProcName: "compiler"},
			},
			Inbound: map[string]*map[string][]workflow.Inbound{
				"testcase": {
					"input": {{
						Name:       "run",
						LabelIndex: 1, // stdin
					}},
					"answer": {{
						Name:       "check",
						LabelIndex: 1, // ans
					}},
				},
				"option": {
					"limitation": {{
						Name:       "run",
						LabelIndex: 2, // limit
					}},
					"compilescript": {{
						Name:       "compile",
						LabelIndex: 1, // script
					}},
				},
				"submission": {
					"source": {{
						Name:       "compile",
						LabelIndex: 0, // source
					}},
				},
			},
		},
		Analyzer: workflow.DefaultAnalyzer{},
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
