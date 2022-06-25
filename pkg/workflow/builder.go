package workflow

import (
	"fmt"

	"github.com/sshwy/yaoj-core/pkg/processor"
)

// Build a workflow
type Builder struct {
	node          map[string]Node
	inbound, edge [][4]string
	err           error
}

func (r *Builder) tryInit() {
	if r.node == nil {
		r.node = map[string]Node{}
	}
	if r.edge == nil {
		r.edge = [][4]string{}
	}
	if r.inbound == nil {
		r.inbound = [][4]string{}
	}
}
func (r *Builder) SetNode(name string, procName string, key bool) {
	r.tryInit()
	r.node[name] = Node{
		ProcName: procName,
		Key:      key,
	}
}

func (r *Builder) AddEdge(from, frlabel, to, tolabel string) {
	r.tryInit()
	r.edge = append(r.edge, [4]string{from, frlabel, to, tolabel})
}

type Groupname string

const (
	Gtests  Groupname = "tests"
	Gsubt   Groupname = "Subtask"
	Gstatic Groupname = "static"
	Gsubm   Groupname = "submission"
)

func (r *Builder) AddInbound(group Groupname, field, to, tolabel string) {
	r.tryInit()
	if group != Gtests && group != Gstatic && group != Gsubm && group != Gsubt {
		r.err = fmt.Errorf("invalid group %s", group)
		return
	}
	r.inbound = append(r.inbound, [4]string{string(group), field, to, tolabel})
}

func (r *Builder) WorkflowGraph() (*WorkflowGraph, error) {
	if r.err != nil {
		return nil, r.err
	}
	graph := NewGraph()
	for name, node := range r.node {
		graph.Node[name] = node
	}
	for _, edge := range r.edge {
		from, frlabel := edge[0], edge[1]
		if _, ok := graph.Node[from]; !ok {
			return nil, fmt.Errorf("invalid edge %v", edge)
		}
		to, tolabel := edge[2], edge[3]
		if _, ok := graph.Node[to]; !ok {
			return nil, fmt.Errorf("invalid edge %v", edge)
		}
		fout := processor.OutputLabel(graph.Node[from].ProcName)
		a := findIndex(fout, frlabel)
		tin := processor.InputLabel(graph.Node[to].ProcName)
		b := findIndex(tin, tolabel)
		if a == -1 || b == -1 {
			return nil, fmt.Errorf("invalid edge %v", edge)
		}
		graph.Edge = append(graph.Edge, Edge{
			From: Outbound{Name: from, LabelIndex: a},
			To:   Inbound{Name: to, LabelIndex: b},
		})
	}
	for _, edge := range r.inbound {
		group, field := edge[0], edge[1]
		to, tolabel := edge[2], edge[3]
		if _, ok := graph.Node[to]; !ok {
			return nil, fmt.Errorf("invalid edge %v", edge)
		}
		tin := processor.InputLabel(graph.Node[to].ProcName)
		b := findIndex(tin, tolabel)
		if b == -1 {
			return nil, fmt.Errorf("invalid edge %v", edge)
		}
		if graph.Inbound[group] == nil {
			graph.Inbound[group] = &map[string][]Inbound{}
		}
		grp := *graph.Inbound[group]
		if grp[field] == nil {
			grp[field] = []Inbound{}
		}
		grp[field] = append(grp[field], Inbound{
			Name:       to,
			LabelIndex: b,
		})
	}
	return &graph, nil
}

func findIndex(s []string, t string) int {
	for i, str := range s {
		if str == t {
			return i
		}
	}
	return -1
}
