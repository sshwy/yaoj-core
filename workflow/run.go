package workflow

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/sshwy/yaoj-core/internal/judger"
	"github.com/sshwy/yaoj-core/utils"
)

var logger = log.New(os.Stderr, "[workflow] ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix)

type RuntimeNode struct {
	Node
	// paths of input files
	Input []string
	// paths of output files
	Output []string
	// result of processor
	Result *judger.Result
}

func (r *RuntimeNode) inputFullfilled() bool {
	for _, path := range r.Input {
		if path == "" {
			return false
		}
	}
	return true
}

// perform a workflow in a directory.
// inboundPath: map[datagroup_name]*map[field]filename
func Run(w Workflow, dir string, inboundPath map[string]*map[string]string, fullscore float64) (*Result, error) {
	if err := w.Valid(); err != nil {
		return nil, fmt.Errorf("workflow validation: %s", err.Error())
	}
	nodes := utils.Map(w.Node, func(node Node) RuntimeNode {
		inLabel, ouLabel := node.Processor().Label()
		return RuntimeNode{
			Node:   node,
			Input:  make([]string, len(inLabel)),
			Output: make([]string, len(ouLabel)),
		}
	})
	if len(w.Inbound) != len(inboundPath) {
		return nil, fmt.Errorf("invalid inboundPath")
	}
	for i, group := range w.Inbound {
		if len(*w.Inbound[i]) != len(*inboundPath[i]) {
			return nil, fmt.Errorf("invalid inboundPath")
		}
		for j, bounds := range *group {
			for _, bound := range bounds {
				nodes[bound.Index].Input[bound.LabelIndex] = (*inboundPath[i])[j]
			}
		}
	}

	for id, node := range nodes {
		if !node.inputFullfilled() {
			panic(fmt.Errorf("input not fullfilled"))
		}
		for i := 0; i < len(node.Output); i++ {
			node.Output[i] = path.Join(dir, utils.RandomString(10))
		}
		for _, edge := range w.EdgeFrom(id) {
			nodes[edge.To.Index].Input[edge.To.LabelIndex] = nodes[edge.From.Index].Output[edge.From.LabelIndex]
		}
		logger.Printf("run node[%d]: input %+v output %+v", id, node.Input, node.Output)
		result := node.Processor().Run(node.Input, node.Output)
		nodes[id].Result = result
	}

	res := w.Analyze(nodes, fullscore)
	return &res, nil
}
