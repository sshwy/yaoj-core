package workflow

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/sshwy/yaoj-core/judger"
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
		return RuntimeNode{
			Node:   node,
			Input:  make([]string, len(node.InEdge)),
			Output: make([]string, len(node.OutEdge)),
		}
	})
	if len(w.Inbound) != len(inboundPath) {
		return nil, fmt.Errorf("invalid inboundPath")
	}
	for i, group := range w.Inbound {
		if len(*w.Inbound[i]) != len(*inboundPath[i]) {
			return nil, fmt.Errorf("invalid inboundPath")
		}
		for j, bound := range *group {
			nodes[bound.Index].Input[bound.LabelIndex] = (*inboundPath[i])[j]
		}
	}

	for id, node := range nodes {
		if !node.inputFullfilled() {
			panic(fmt.Errorf("input not fullfilled"))
		}
		for i, edge := range node.OutEdge {
			node.Output[i] = path.Join(dir, utils.RandomString(10))
			if edge.Bound != nil {
				nodes[edge.Bound.Index].Input[edge.Bound.LabelIndex] = node.Output[i]
			}
		}
		logger.Printf("run node[%d]: input %+v output %+v", id, node.Input, node.Output)
		result, err := node.Processor().Run(node.Input, node.Output)
		if err != nil {
			return nil, err
		}
		nodes[id].Result = result
	}

	res := w.Analyze(nodes, fullscore)
	return &res, nil
}
