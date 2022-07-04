package workflow

import (
	"fmt"
	"log"
	goPlugin "plugin"

	"github.com/sshwy/yaoj-core/pkg/processor"
	"github.com/sshwy/yaoj-core/pkg/utils"
)

// Analyzer generates result of a workflow.
type Analyzer interface {
	Analyze(w Workflow, nodes map[string]RuntimeNode, fullscore float64) Result
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

func (r DefaultAnalyzer) Analyze(w Workflow, nodes map[string]RuntimeNode, fullscore float64) Result {
	for _, bounds := range *w.Inbound[Gsubm] {
		for _, bound := range bounds {
			nodes[bound.Name].Attr["dependon"] = "user"
		}
	}
	for {
		flag := false
		for _, edge := range w.Edge {
			if nodes[edge.From.Name].Attr["dependon"] == "user" &&
				nodes[edge.To.Name].Attr["dependon"] != "user" {

				nodes[edge.To.Name].Attr["dependon"] = "user"
				flag = true
			}
		}
		if !flag {
			break
		}
	}
	res := Result{
		ResultMeta: ResultMeta{
			Score:     fullscore,
			Fullscore: fullscore,
			Title:     "Accepted",
		},
		File: []ResultFileDisplay{},
	}

	for name, node := range nodes {
		if node.Result == nil {
			continue
		}
		log.Printf("node[%s]: %d", name, node.Result.Code)
		labels := processor.OutputLabel(node.ProcName)
		list := []ResultFileDisplay{}
		for i, label := range labels {
			list = append(list, FileDisplay(node.Output[i], label, 5000))
		}
		list = append(list, ResultFileDisplay{
			Title:   "message",
			Content: node.Result.Msg,
		})
		if node.Result.Code != processor.Ok {
			if node.Attr["dependon"] == "user" {
				var title = "Unaccepted"
				switch node.Result.Code {
				case processor.TimeExceed:
					title = "Time Limit Exceed"
				case processor.RuntimeError:
					title = "Runtime Error"
				case processor.DangerousSyscall:
					title = "Dangerous System Call"
				case processor.ExitError:
					title = "Exit Code Error"
				case processor.OutputExceed:
					title = "Output Limit Exceed"
				case processor.MemoryExceed:
					title = "Memory Limit Exceed"
				}

				return Result{
					ResultMeta: ResultMeta{
						Score:     0,
						Fullscore: fullscore,
						Title:     title,
					},
					File: list,
				}
			} else { // system error
				return Result{
					ResultMeta: ResultMeta{
						Score:     0,
						Fullscore: fullscore,
						Title:     "System Error",
					},
					File: list,
				}
			}
		}
		if node.Key {
			res.ResultMeta.Memory += utils.ByteValue(*node.Result.Memory)
			res.ResultMeta.Time += *node.Result.CpuTime
			res.File = append(res.File, list...)
		}
	}
	return res
}

var _ Analyzer = DefaultAnalyzer{}
