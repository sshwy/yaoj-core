package workflow

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	goPlugin "plugin"

	"github.com/k0kubun/pp/v3"
	"github.com/sshwy/yaoj-core/pkg/processor"
	"github.com/sshwy/yaoj-core/pkg/utils"
	"golang.org/x/text/encoding/charmap"
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
		labels := processor.OutputLabel(node.ProcName)
		list := []ResultFileDisplay{}
		for i, label := range labels {
			list = append(list, FileDisplay(node.Output[i], label, 5000))
		}
		list = append(list, ResultFileDisplay{
			Title:   "message",
			Content: name + ": " + node.Result.Msg,
		})

		if node.Key {
			res.ResultMeta.Memory += utils.ByteValue(*node.Result.Memory)
			res.ResultMeta.Time += *node.Result.CpuTime
			res.File = append(res.File, list...)
		}
	}

	for name, node := range nodes {
		if node.Result == nil {
			continue
		}
		if node.Result.Code != processor.Ok {
			labels := processor.OutputLabel(node.ProcName)
			list := []ResultFileDisplay{}
			if !node.Key { // 如果是关键结点那么文件啥的已经被展示了
				for i, label := range labels {
					list = append(list, FileDisplay(node.Output[i], label, 5000))
				}
				list = append(list, ResultFileDisplay{
					Title:   "message",
					Content: name + ": " + node.Result.Msg,
				})
			}

			if node.ProcName == "checker:testlib" {
				type Result struct {
					XMLName xml.Name `xml:"result"`
					Outcome string   `xml:"outcome,attr"`
				}
				var result Result
				file, _ := os.Open(node.Output[0])
				// parse xml encoded windows1251
				d := xml.NewDecoder(file)
				d.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
					switch charset {
					case "windows-1251":
						return charmap.Windows1251.NewDecoder().Reader(input), nil
					default:
						return nil, fmt.Errorf("unknown charset: %s", charset)
					}
				}
				d.Decode(&result)

				res.File = append(res.File, list...)
				res.File = append(res.File, FileDisplay(node.Input[3], "answer", 5000))
				if result.Outcome != "accepted" {
					res.Title = "Wrong Answer"
					res.Score = 0
					pp.Print(result)
				}
				return res
			}
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

				res.Title = title
				res.Score = 0
				res.File = append(res.File, list...)
				return res
			}
			// system error
			res.Title = "System Error"
			res.Score = 0
			res.File = append(res.File, list...)
			return res
		}
	}
	return res
}

var _ Analyzer = DefaultAnalyzer{}
