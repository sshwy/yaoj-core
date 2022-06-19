/*
Package problem is a package for test problem manipulation.

Overview

A test problem is store in a directory, for example, "dir", with following
structure:

  |-dir/problem.json
  |-dir/data/
  |---dir/tests/
  |---dir/subtasks/
  |---dir/static/
  *-dir/statement/ stores problem statement
  *---dir/statement/statement.md stores markdown statement
  |-dir/workflow/ stores test workflow
  |---dir/workflow/graph.json stores the workflow graph
  |-dir/patch/ stores added files

Statement

Markdown is the standard format of statement. Other formats may be supported in
the future.

Workflow

See github.com/sshwy/yaoj-core/workflow.

Problem gives workflow 4 datagroups naming "testcase", "subtask", "static" and
"submission" respectively.
*/
package problem
