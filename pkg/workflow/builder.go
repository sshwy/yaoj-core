package workflow

import "fmt"

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

func (r *Builder) AddInbound(group, field, to, tolabel string) {
	r.tryInit()
	if group != "testcase" && group != "subtask" && group != "submission" && group != "static" {
		r.err = fmt.Errorf("invalid group %s", group)
		return
	}
	r.inbound = append(r.inbound, [4]string{group, field, to, tolabel})
}

func (r *Builder) WorkflowGraph() (*WorkflowGraph, error) {
	if r.err != nil {
		return nil, r.err
	}
	graph := NewGraph()
	for name, node := range r.node {
		if node.Processor() == nil {
			return nil, fmt.Errorf("invalid processor %s", node.ProcName)
		}
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
		_, fout := graph.Node[from].Processor().Label()
		a := findIndex(fout, frlabel)
		tin, _ := graph.Node[to].Processor().Label()
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
		tin, _ := graph.Node[to].Processor().Label()
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