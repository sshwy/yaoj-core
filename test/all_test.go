package test_test

import "testing"

func TestAll(t *testing.T) {
	probDataDir = t.TempDir()
	problemDumpDir = t.TempDir()
	t.Run("MakeWorkflowGraph", MakeWorkflowGraph)
	t.Run("MakeProbData", MakeProbData)
	t.Run("LoadProblem", LoadProblem)
	t.Run("DumpProblem", DumpProblem)
	t.Run("ExtractProblem", ExtractProblem)
}
