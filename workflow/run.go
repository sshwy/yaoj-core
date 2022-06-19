package workflow

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/k0kubun/pp/v3"
	"github.com/sshwy/yaoj-core/internal/judger"
	"github.com/sshwy/yaoj-core/utils"
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
	Result *judger.Result
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
	hash := sha256.New()
	hash.Write(r.hash[:])
	res = make([]sha, len(r.Output))
	_, labels := r.Processor().Label()
	for i, label := range labels {
		hash.Write([]byte(label))
		res[i] = *(*sha)(hash.Sum(nil))
	}
	return
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
		return nil, fmt.Errorf("invalid inboundPath: missing field")
	}
	for i, group := range w.Inbound {
		for j, bounds := range *group {
			if _, ok := (*inboundPath[i])[j]; !ok {
				return nil, fmt.Errorf("invalid inboundPath: missing field %s %s", i, j)
			}
			for _, bound := range bounds {
				nodes[bound.Index].Input[bound.LabelIndex] = (*inboundPath[i])[j]
			}
		}
	}

	for id, node := range nodes {
		if !node.inputFullfilled() {
			panic(fmt.Errorf("input not fullfilled"))
		}
		// node.calcHash()
		// log.Print(node.outputHash())
		// log.Printf("%d, %v", id, node.hash)
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
