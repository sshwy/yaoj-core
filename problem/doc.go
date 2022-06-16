/*
Package problem is a package for test problem manipulation.

Overview

A test problem is store in a directory, for example, "dir", with following
structure:

  |-dir/datagroup/ stores all datagroups (e. g. testcase, compileoption, checker, std)
  |---dir/datagroup/xxx/ stores a datagroup
  |---dir/datagroup/submission/ stores a datagroup denoting submit format (special)
  |-dir/statement/ stores problem statement
  |---dir/statement/statement.md stores markdown statement
  |-dir/workflow/ stores test workflow
  |---dir/workflow/graph.json stores the workflow graph
  |---dir/workflow/analyzer.go stores custom analyzer (go plugin)

Datagroup

Datagroup (ProbDtgp) looks like a table consists of several records, each of
which contains a series of fields and its corresponding value. All records in
one datagroup possess the same fields.

A datagroup is stored in a directory, for example, "dirdgtp", with its records
stored in the following format:

  [record id].[field].[arbtrary suffix]

"record id" is a integer starting from 0.
"field" namely is the field's name.
"arbtrary suffix" namely is an arbtrary suffix, which for human readability.

Note that datagroup directory should not contain sub directories.

Statement

Markdown is the standard format of statement. Other formats may be supported in
the future.

Workflow

See github.com/sshwy/yaoj-core/workflow
*/
package problem
