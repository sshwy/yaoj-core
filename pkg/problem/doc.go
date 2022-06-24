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
  |-dir/workflow/ stores test workflow
  |---dir/workflow/graph.json stores the workflow graph
  |-dir/statement/ stores statement files
  |---dir/statement/s.[lang][id].md stores statement files
  |---dir/statement/t.[lang][id].md stores tutorial files
  |---dir/statement/xxx stores assert files
  |-dir/patch/ stores added files

Tests

Usually a problem contains multiple testcases, whose data are stored in
dir/data/tests/, since a contestant's submission should be tested enoughly
before being considered correct. All testcases have the same fields.

Testcase may have "_score" field, whose value is either a number string or
"average".

Subtask

To better assign points to testcases with different intensity, it's common to
set up several subtasks for the problem, each containing a series of testcases.

Note that if subtask is enabled, independent testcases (i.e. not in a subtask)
are not allowed.

If subtask data occured (at least one field, at least one record), subtask is
enabled.

For some problem, different subtasks use different files to test contestant'
submission, for example, checker or input generator. Thus these data are stored
in dir/data/tests/. Again, all subtasks' data have the same fields.

Subtask may have "_score" field, whose value is a number string.

For problem enabling subtask, its testcase and subtask both contain
"_subtaskid" field, determining which subtask the testcase belongs to.

Static Data

Common data are stored in dir/data/static/ shared by all testcases.

Testcase Score

Score of a testcase is calculated as follows:

If subtask is enabled, testcase's "_score" is ignored, its score is
{subtask score} / {number of tests in this subtask}.

Otherwise, if "average" is specified for "_score" field, its score is
{problem total score} / {number of testcases}. Else its "_score" should be
a number string, denoting its score.

Statement

Markdown is the standard format of statement. Other formats may be supported in
the future.

Workflow

See github.com/sshwy/yaoj-core/workflow.

Problem gives workflow 4 datagroups naming "testcase", "subtask", "static" and
"submission" respectively.
*/
package problem
