package workflow

import (
	"os"
	"strings"
	"time"

	"github.com/sshwy/yaoj-core/processor"
	"github.com/sshwy/yaoj-core/utils"
)

type Bound struct {
	// index of the node in the array
	Index int
	// index of the file in input (output) array
	LabelIndex int
}

// When it comes to out edge, label denotes output label and bound denotes the destination.
// Otherwise (in edge), label denotes input label and bound denotes the origin.
// Actually an edge is just a file in os.
// If a Edge has no Bound, it should be workflow inbound(outbound) edge.
type Edge struct {
	Bound *Bound
	Label string
}

type Node struct {
	// processor name
	ProcName string
	InEdge   []Edge
	OutEdge  []Edge
}

func (r *Node) Processor() processor.Processor {
	return processor.Get(r.ProcName)
}

// workflow describes how to perform a single testcase's judgement
type Workflow struct {
	// a node itself is just a processor
	Node []Node
	// inbound consists a series of data group
	Inbound []DataBoundGroup
}

// for storage
// func (r *Workflow) Serialize() []byte

// check whether it's a well-formatted DAG, its inbound coverage and sth else
func (r *Workflow) Valid() error {
	logger.Printf("Workflow.Valid not implemented!!!")
	return nil
}

// transform to dot file content
// func (r *Workflow) Dot() string

// parse dot file to workflow
// func (r *Workflow) ParseDot(content string) error

// use a string represents a data field
type DataLabel string

// connect data with bound
type DataBound struct {
	Data  DataLabel
	Bound Bound
}

// a series of data
type DataBoundGroup []DataBound

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

// Analyzer generates result of a workflow.
type Analyzer interface {
	Analyze(nodes []RuntimeNode, fullscore float64) Result
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
