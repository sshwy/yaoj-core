package main

import (
	"github.com/sshwy/yaoj-core/internal/judger"
	"github.com/sshwy/yaoj-core/workflow"
)

type analyzerPlugin struct {
}

func (r analyzerPlugin) Analyze(nodes []workflow.RuntimeNode, fullscore float64) workflow.Result {
	for _, node := range nodes {
		if node.Result.Code != judger.Ok {
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
