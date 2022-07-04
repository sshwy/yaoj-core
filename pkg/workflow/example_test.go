package workflow_test

import (
	"log"

	"github.com/sshwy/yaoj-core/pkg/workflow"
)

func ExampleBuilder() {
	var b workflow.Builder
	b.SetNode("compile", "compiler", false)
	b.SetNode("run", "runner:stdio", true)
	b.SetNode("check", "checker:hcmp", false)
	b.AddInbound(workflow.Gsubm, "source", "compile", "source")
	b.AddInbound(workflow.Gstatic, "compilescript", "compile", "script")
	b.AddInbound(workflow.Gstatic, "limitation", "run", "limit")
	b.AddInbound(workflow.Gtests, "input", "run", "stdin")
	b.AddInbound(workflow.Gtests, "answer", "check", "ans")
	b.AddEdge("compile", "result", "run", "executable")
	b.AddEdge("run", "stdout", "check", "out")
	graph, err := b.WorkflowGraph()
	if err != nil {
		log.Print(err)
	}
	log.Print(string(graph.Serialize()))
}
