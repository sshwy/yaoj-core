package workflow

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/k0kubun/pp/v3"
	"github.com/sshwy/yaoj-core/pkg/private/processors"
	"github.com/sshwy/yaoj-core/pkg/processor"
	"github.com/sshwy/yaoj-core/pkg/utils"
)

var logger = log.New(os.Stderr, "[workflow] ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix)

type sha [32]byte

func (r sha) String() string {
	s := ""
	for _, v := range r {
		s += fmt.Sprintf("%02x", v)
	}
	return s
}

type RuntimeNode struct {
	Node
	// paths of input files
	Input []string
	// paths of output files
	Output []string
	// result of processor
	Result *processor.Result
	hash   sha
}

func (r *RuntimeNode) inputFullfilled() bool {
	for _, path := range r.Input {
		if path == "" {
			return false
		}
	}
	return true
}

// SHA256 hash for file content.
// for any error, return empty hash
func fileHash(name string) sha {
	hash := sha256.New()
	f, err := os.Open(name)
	if err != nil {
		return sha{}
	}
	defer f.Close()

	if _, err := io.Copy(hash, f); err != nil {
		return sha{}
	}
	var b = hash.Sum(nil)
	// pp.Print(b)
	if len(b) != 32 {
		pp.Print(b)
		panic(b)
	}
	return *(*sha)(b)
}

// sum up hash of all input files
func (r *RuntimeNode) calcHash() {
	hash := sha256.New()
	for _, path := range r.Input {
		hashval := fileHash(path)
		log.Print(path, " ", hashval)
		hash.Write(hashval[:])
	}
	var b = hash.Sum(nil)
	// pp.Print(b)
	if len(b) != 32 {
		pp.Print(b)
		panic(b)
	}
	r.hash = *(*sha)(b)
}

// generate hash for output files
func (r *RuntimeNode) outputHash() (res []sha) {
	if r.hash == (sha{}) {
		r.calcHash()
	}
	_, labels := r.Processor().Label()
	res = make([]sha, len(r.Output))
	hash := sha256.New()
	for i, label := range labels {
		hash.Reset()
		hash.Write(r.hash[:])
		hash.Write([]byte(label))
		res[i] = *(*sha)(hash.Sum(nil))
	}
	return
}

func runtimeNodes(node map[string]Node) (res map[string]RuntimeNode) {
	res = map[string]RuntimeNode{}
	for k, v := range node {
		inLabel, ouLabel := processors.Get(v.ProcName).Label()
		res[k] = RuntimeNode{
			Node:   v,
			Input:  make([]string, len(inLabel)),
			Output: make([]string, len(ouLabel)),
		}
	}
	return
}

func topologicalEnum(w Workflow, handler func(id string) error) error {
	indegree := map[string]int{}
	for _, edge := range w.Edge {
		indegree[edge.To.Name]++
	}
	for {
		flag := false
		p := ""
		for id := range w.Node {
			if indegree[id] == 0 {
				flag = true
				p = id
				break
			}
		}
		if !flag {
			break
		}
		logger.Printf("topo current id=%s", p)
		indegree[p] = -1
		for _, edge := range w.EdgeFrom(p) {
			indegree[edge.To.Name]--
		}
		err := handler(p)
		if err != nil {
			return err
		}
	}
	for id := range w.Node {
		if indegree[id] != -1 {
			return fmt.Errorf("invalid DAG! id=%s", id)
		}
	}
	return nil
}

// Get the processor of the node.
func (r RuntimeNode) Processor() processor.Processor {
	return processors.Get(r.ProcName)
}

// perform a workflow in a directory.
// inboundPath: map[datagroup_name]*map[field]filename
func Run(w Workflow, dir string, inboundPath map[string]*map[string]string, fullscore float64) (*Result, error) {
	nodes := runtimeNodes(w.Node)

	if len(w.Inbound) != len(inboundPath) {
		return nil, fmt.Errorf("invalid inboundPath: missing field")
	}
	for i, group := range w.Inbound {
		if group == nil {
			return nil, fmt.Errorf("w.Inbound[%s] == nil", i)
		}
		data := inboundPath[i]
		if data == nil {
			return nil, fmt.Errorf("inboundPath[%s] == nil", i)
		}
		for j, bounds := range *group {
			if _, ok := (*data)[j]; !ok {
				return nil, fmt.Errorf("invalid inboundPath: missing field %s %s", i, j)
			}
			for _, bound := range bounds {
				nodes[bound.Name].Input[bound.LabelIndex] = (*inboundPath[i])[j]
			}
		}
	}

	err := topologicalEnum(w, func(id string) error {
		node := nodes[id]
		if !node.inputFullfilled() {
			return fmt.Errorf("input not fullfilled")
		}
		node.calcHash()
		// log.Print(node.outputHash())
		// log.Printf("%d, %v", id, node.hash)
		for i := 0; i < len(node.Output); i++ {
			node.Output[i] = path.Join(dir, utils.RandomString(10))
		}
		for _, edge := range w.EdgeFrom(id) {
			nodes[edge.To.Name].Input[edge.To.LabelIndex] = nodes[edge.From.Name].Output[edge.From.LabelIndex]
		}
		logger.Printf("run node[%s]:", id)
		logger.Printf("input %+v", node.Input)
		logger.Printf("output %+v", node.Output)
		result := node.Processor().Run(node.Input, node.Output)
		nd := nodes[id]
		nd.Result = result
		nodes[id] = nd
		return nil
	})
	if err != nil {
		return nil, err
	}

	res := w.Analyze(nodes, fullscore)
	return &res, nil
}
