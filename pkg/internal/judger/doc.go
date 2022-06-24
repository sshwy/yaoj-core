/*
Package judger is the go wrapper for https://github.com/sshwy/yaoj-judger.

Since this package requires a lot of depedencies and a "go generate" to
initialize, it's not recommended to import this package.

To build (run) this package, you need to make sure auditd, flex, make,
gengetopt, bison, xxd, strace, and clang toolkit (basically clang++) is
available via command line. If not, install them.  Before building, run go
generate for some necessary files.
*/
package judger
