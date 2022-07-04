package main

import (
	"github.com/sshwy/yaoj-core/pkg/processor"
	"github.com/sshwy/yaoj-core/pkg/workflow"
)

type analyzerPlugin struct {
}

func (r analyzerPlugin) Analyze(w workflow.Workflow, nodes map[string]workflow.RuntimeNode, fullscore float64) workflow.Result {
	for _, node := range nodes {
		if node.Result.Code != processor.Ok {
			return workflow.Result{
				ResultMeta: workflow.ResultMeta{
					Score:     0,
					Fullscore: fullscore,
				},
			}
		}
	}
	return workflow.Result{
		ResultMeta: workflow.ResultMeta{
			Score:     fullscore,
			Fullscore: fullscore,
		},
	}
}

var AnalyzerPlugin workflow.Analyzer = analyzerPlugin{}
