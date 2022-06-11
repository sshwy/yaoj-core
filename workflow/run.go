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

type runtimeNode struct {
	Node
	input  []string
	output []string
}

func (r *runtimeNode) inputFullfilled() bool {
	for _, path := range r.input {
		if path == "" {
			return false
		}
	}
	return true
}

// perform a workflow in a preconfigured safe environment
func Run(w Workflow, dir string, inboundPath [][]string, a Analyzer) (*Result, error) {
	if err := w.Valid(); err != nil {
		return nil, fmt.Errorf("workflow validation: %s", err.Error())
	}
	nodes := utils.Map(w.Node, func(node Node) runtimeNode {
		return runtimeNode{
			Node:   node,
			input:  make([]string, len(node.InEdge)),
			output: make([]string, len(node.OutEdge)),
		}
	})
	if len(w.Inbound) != len(inboundPath) {
		return nil, fmt.Errorf("invalid inboundPath")
	}
	for i, group := range w.Inbound {
		if len(w.Inbound[i]) != len(inboundPath[i]) {
			return nil, fmt.Errorf("invalid inboundPath")
		}
		for j, bound := range group {
			nodes[bound.Bound.Index].input[bound.Bound.LabelIndex] = inboundPath[i][j]
		}
	}
	results := make([]*judger.Result, len(w.Node))
	for id, node := range nodes {
		if !node.inputFullfilled() {
			panic(fmt.Errorf("input not fullfilled"))
		}
		for i, edge := range node.OutEdge {
			node.output[i] = path.Join(dir, utils.RandomString(10))
			if edge.Bound != nil {
				nodes[edge.Bound.Index].input[edge.Bound.LabelIndex] = node.output[i]
			}
		}
		logger.Printf("run node[%d]: input %v output %v", id, node.input, node.output)
		result, err := node.Processor().Run(node.input, node.output)
		if err != nil {
			return nil, err
		}
		results[id] = result
	}
	outboundPath := make([]string, len(w.Outbound))

	for i, bound := range w.Outbound {
		outboundPath[i] = nodes[bound.Bound.Index].output[bound.Bound.LabelIndex]
	}
	res := a.Analyze(results, outboundPath)
	return &res, nil
}
