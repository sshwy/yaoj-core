package workflow

import (
	"fmt"
	goPlugin "plugin"

	"github.com/sshwy/yaoj-core/pkg/processor"
	"github.com/sshwy/yaoj-core/pkg/utils"
)

// Analyzer generates result of a workflow.
type Analyzer interface {
	Analyze(nodes map[string]RuntimeNode, fullscore float64) Result
}

func LoadAnalyzer(plugin string) (Analyzer, error) {
	p, err := goPlugin.Open(plugin)
	if err != nil {
		return nil, err
	}

	label, err := p.Lookup("AnalyzerPlugin")
	if err != nil {
		return nil, err
	}
	analyzer, ok := label.(*Analyzer)
	if ok {
		return *analyzer, nil
	} else {
		return nil, fmt.Errorf("AnalyzerPlugin not implement Analyzer")
	}
}

type DefaultAnalyzer struct {
}

func (r DefaultAnalyzer) Analyze(nodes map[string]RuntimeNode, fullscore float64) Result {
	res := Result{
		ResultMeta: ResultMeta{
			Score:     fullscore,
			Fullscore: fullscore,
			Title:     "Accepted",
		},
	}

	for i, node := range nodes {
		if node.Result == nil {
			continue
		}
		if node.Result.Code != processor.Ok {
			return Result{
				ResultMeta: ResultMeta{
					Score:     0,
					Fullscore: fullscore,
					Title:     "Not Accepted",
				},
				File: []ResultFileDisplay{
					{
						Title:   "Error Node",
						Content: fmt.Sprintf("name=%s, proc=%s, Code=%v", i, node.ProcName, node.Result.Code),
					},
				},
			}
		}
		if node.Key {
			res.ResultMeta.Memory += utils.ByteValue(*node.Result.Memory)
			res.ResultMeta.Time += *node.Result.CpuTime
		}
	}
	return res
}

var _ Analyzer = DefaultAnalyzer{}
