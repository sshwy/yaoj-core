package workflow_test

import (
	"testing"

	"github.com/sshwy/yaoj-core/judger"
	"github.com/sshwy/yaoj-core/workflow"
)

type testanalyzer struct {
}

func (r testanalyzer) Analyze(results []*judger.Result, outbound []string) workflow.Result {
	if results[2].Code == judger.Ok {
		return workflow.Result{
			Score: 100,
		}
	}
	return workflow.Result{
		Score: 0,
	}
}

func TestWorkflow(t *testing.T) {
	w := workflow.Workflow{
		Node: []workflow.Node{
			{
				ProcName: "compiler",
				InEdge:   []workflow.Edge{{Label: "source"}, {Label: "script"}},
				OutEdge: []workflow.Edge{
					{
						Label: "result",
						Bound: &workflow.Bound{
							Index:      1, // runner
							LabelIndex: 0, // executable
						},
					},
					{Label: "log"},
					{Label: "judgerlog"},
				},
			},
			{
				ProcName: "runner:stdio",
				InEdge: []workflow.Edge{
					{
						Label: "executable",
						Bound: &workflow.Bound{
							Index:      0, // compiler
							LabelIndex: 0, // result
						},
					},
					{Label: "stdin"},
					{Label: "limit"},
				},
				OutEdge: []workflow.Edge{
					{
						Label: "stdout",
						Bound: &workflow.Bound{
							Index:      2, // checker
							LabelIndex: 0, // out,
						},
					},
					{Label: "stderr"},
					{Label: "judgerlog"},
				},
			},
			{
				ProcName: "checker:hcmp",
				InEdge: []workflow.Edge{
					{
						Label: "out",
						Bound: &workflow.Bound{
							Index:      1, // runner
							LabelIndex: 0, // stdout
						},
					},
					{Label: "ans"},
				},
				OutEdge: []workflow.Edge{{Label: "result"}},
			},
		},
		Inbound: []workflow.DataBoundGroup{
			[]workflow.DataBound{
				{
					Data: "input",
					Bound: workflow.Bound{
						Index:      1, // runner:stdio
						LabelIndex: 1, // stdin
					},
				},
				{
					Data: "answer",
					Bound: workflow.Bound{
						Index:      2, // checker
						LabelIndex: 1, // ans
					},
				},
			},
			[]workflow.DataBound{
				{
					Data: "limitation",
					Bound: workflow.Bound{
						Index:      1, // runner:stdio
						LabelIndex: 2, // limit
					},
				},
			},
			[]workflow.DataBound{
				{
					Data: "compilescript",
					Bound: workflow.Bound{
						Index:      0, // compiler
						LabelIndex: 1, // script
					},
				},
			},
			[]workflow.DataBound{
				{
					Data: "source",
					Bound: workflow.Bound{
						Index:      0, // compiler
						LabelIndex: 0, // source
					},
				},
			},
		},
		Outbound: []workflow.DataBound{},
	}
	dir := t.TempDir()
	res, err := workflow.Run(w, dir, [][]string{
		{"testdata/main.in", "testdata/main.ans"},
		{"testdata/main.lim"},
		{"testdata/script.sh"},
		{"testdata/main.cpp"},
	}, testanalyzer{})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("score: ", res.Score)
}
