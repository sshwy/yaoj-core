package workflow_test

import (
	"testing"

	"github.com/k0kubun/pp/v3"
	"github.com/sshwy/yaoj-core/judger"
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
	}
	dir := t.TempDir()
	res, err := workflow.Run(w, dir, [][]string{
		{"testdata/main.in", "testdata/main.ans"},
		{"testdata/main.lim"},
		{"testdata/script.sh"},
		{"testdata/main.cpp"},
	}, testanalyzer{}, 100)
	if err != nil {
		t.Error(err)
		return
	}
	if res.Score != res.Fullscore {
		t.Errorf("score=%f, expect %f", res.Score, res.Fullscore)
		return
	}
	t.Log(pp.Sprint(*res))
}
