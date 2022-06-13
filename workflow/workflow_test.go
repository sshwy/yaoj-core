package workflow_test

import (
	"testing"

	"github.com/awalterschulze/gographviz"
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
		Inbound: map[string]*map[string]workflow.Bound{
			"testcase": {
				"input": {
					Index:      1, // runner:stdio
					LabelIndex: 1, // stdin
				},
				"answer": {
					Index:      2, // checker
					LabelIndex: 1, // ans
				},
			},
			"option": {
				"limitation": {
					Index:      1, // runner:stdio
					LabelIndex: 2, // limit
				},
				"compilescript": {
					Index:      0, // compiler
					LabelIndex: 1, // script
				},
			},
			"submission": {
				"source": {
					Index:      0, // compiler
					LabelIndex: 0, // source
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
}

//go:generate go build -buildmode=plugin -o ./testdata ./testdata/custom_analyzer.go
func TestLoadAnalyzer(t *testing.T) {
	a, err := workflow.LoadAnalyzer("testdata/custom_analyzer.so")
	if err != nil {
		t.Error(err)
		return
	}
	pp.Print(a)
}

var dotxt string = `
digraph{
  graph[ranksep=1,rankdir=LR,label="Workflow",nodesep=1];
  node[shape=record,style=filled,fillcolor=white];
  edge[arrowhead=odot, color="#9d9d9d"];

  idg for inbound data group
  subgraph cluster_idg_group1{
    graph[label="Data Group 1"];
    node[shape=record,style=filled,fillcolor=grey];
    main;
  } 

  subgraph cluster_idg_group2{
    graph[label="Data Group 2"];
    node[shape=record,style=filled,fillcolor=grey];
    inputopt;
		ans;
  } 
  subgraph processor{
    node[shape=hexagon];
    compile;
		judge;
		check;
		makeinput;
  }

  main -> compile[headlabel="source"];
  compile -> compile_log[taillabel="log"];
  compile -> judge[taillabel="result", headlabel="exec"];
  inputopt -> makeinput[headlabel="option"];
  makeinput -> judge[taillabel="result", headlabel="input"];
  judge -> judger_err[taillabel="mainlog"];
  judge -> judge_log[taillabel="log"];
  judge -> check[taillabel="output", headlabel="out"];
  ans -> check[headlabel="ans"];
  check -> check_log[taillabel="log"];
  check -> score[taillabel="score"];
}
`

func TestDot(t *testing.T) {
	graphAst, err := gographviz.ParseString(dotxt)
	if err != nil {
		t.Error(err)
		return
	}
	graph := gographviz.NewGraph()
	if err := gographviz.Analyse(graphAst, graph); err != nil {
		panic(err)
	}
	// pp.Print(graphAst)
	// pp.Print(graph)
}
