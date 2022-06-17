package workflow

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sshwy/yaoj-core/processor"
	"github.com/sshwy/yaoj-core/utils"
)

type DtgpBound struct {
	Group string
	Field string
}

type Bound struct {
	// index of the node in the array
	Index int
	// index of the file in input (output) array
	LabelIndex int
}

type Inbound Bound

type Outbound Bound

// Edge between nodes.
type Edge struct {
	From Outbound
	To   Inbound
}

type Node struct {
	// processor name
	ProcName string
}

func (r *Node) Processor() processor.Processor {
	return processor.Get(r.ProcName)
}

type WorkflowGraph struct {
	// a node itself is just a processor
	Node []Node
	Edge []Edge
	// inbound consists a series of data group.
	// Inbound: map[datagroup_name]*map[field]Bound
	Inbound map[string]*map[string][]Inbound
}

func (r *WorkflowGraph) Serialize() []byte {
	res, err := json.Marshal(*r)
	if err != nil {
		panic(err)
	}
	return res
}

func (r *WorkflowGraph) EdgeFrom(nodeid int) []Edge {
	res := []Edge{}
	for _, edge := range r.Edge {
		if edge.From.Index == nodeid {
			res = append(res, edge)
		}
	}
	return res
}
func (r *WorkflowGraph) EdgeTo(nodeid int) []Edge {
	res := []Edge{}
	for _, edge := range r.Edge {
		if edge.To.Index == nodeid {
			res = append(res, edge)
		}
	}
	return res
}

// check whether it's a well-formatted DAG, its inbound coverage and sth else
func (r *WorkflowGraph) Valid() error {
	for i, node := range r.Node {
		proc := node.Processor()
		if proc == nil {
			return fmt.Errorf("node[%d] has invalid processor name (%s)", i, node.ProcName)
		}
		for _, edge := range r.EdgeTo(i) {
			if edge.From.Index >= i || edge.From.Index < 0 {
				return fmt.Errorf("invalid Edge %+v", edge)
			}
		}
		for _, edge := range r.EdgeFrom(i) {
			if edge.To.Index <= i || edge.To.Index >= len(r.Node) {
				return fmt.Errorf("invalid Edge %+v", edge)
			}
		}
	}
	for i, group := range r.Inbound {
		for j, bounds := range *group {
			for _, bound := range bounds {
				node := r.Node[bound.Index]
				inLabel, _ := node.Processor().Label()
				if bound.Index >= len(r.Node) || bound.Index < 0 {
					return fmt.Errorf("inbound[%s][%s] has invalid node index %d", i, j, bound.Index)
				}
				if bound.LabelIndex >= len(inLabel) {
					return fmt.Errorf("inbound[%s][%s] has invalid node label index %d", i, j, bound.LabelIndex)
				}
			}
		}
	}
	return nil
}

func Load(serial []byte) (*WorkflowGraph, error) {
	var graph WorkflowGraph
	err := json.Unmarshal(serial, &graph)
	if err != nil {
		return nil, err
	}
	if err := graph.Valid(); err != nil {
		return nil, err
	}
	return &graph, nil
}

func LoadFile(path string) (*WorkflowGraph, error) {
	serial, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Load(serial)
}

// workflow describes how to perform a single testcase's judgement
type Workflow struct {
	*WorkflowGraph
	Analyzer
}

type Result struct {
	Score     float64
	Fullscore float64
	Time      time.Duration
	Memory    utils.ByteValue
	// e. g. "Accepted", "Wrong Answer"
	Title string
	// a list of file content to display
	File []ResultFileDisplay
	// other tags
	Property map[string]string
}

type ResultFileDisplay struct {
	Title   string
	Content string
}

func fetchFileContent(path string, len int) []byte {
	file, err := os.Open(path)
	if err != nil {
		return []byte("[error] " + err.Error())
	}
	defer file.Close()
	b := make([]byte, len)
	file.Read(b)
	return b
}

func FileDisplay(path string, title string, len int) ResultFileDisplay {
	content := strings.TrimRight(string(fetchFileContent(path, len)), "\x00 \n\t\r")
	return ResultFileDisplay{
		Title:   title,
		Content: content,
	}
}
