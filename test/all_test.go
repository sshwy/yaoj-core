package test_test

import "testing"

func TestAll(t *testing.T) {
	probDataDir = t.TempDir()
	problemDumpDir = t.TempDir()
	if !t.Run("MakeWorkflowGraph", MakeWorkflowGraph) {
		return
	}
	t.Run("MakeProbData", MakeProbData)
	t.Run("LoadProblem", LoadProblem)
	t.Run("DumpProblem", DumpProblem)
	t.Run("ExtractProblem", ExtractProblem)
	t.Run("RunProblem", RunProblem)
}
