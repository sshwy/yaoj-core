
package processor

// generated by scripts/genprocs
func init() {
	inLabel[`checker:hcmp`]=[]string{`out`,`ans`}
	ouLabel[`checker:hcmp`]=[]string{`result`}
	inLabel[`checker:testlib`]=[]string{`checker`,`input`,`output`,`answer`}
	ouLabel[`checker:testlib`]=[]string{`xmlreport`,`stderr`,`judgerlog`}
	inLabel[`compiler`]=[]string{`source`,`script`}
	ouLabel[`compiler`]=[]string{`result`,`log`,`judgerlog`}
	inLabel[`generator:testlib`]=[]string{`generator`,`arguments`}
	ouLabel[`generator:testlib`]=[]string{`output`,`stderr`,`judgerlog`}
	inLabel[`inputmaker`]=[]string{`source`,`option`,`generator`}
	ouLabel[`inputmaker`]=[]string{`result`,`stderr`,`judgerlog`}
	inLabel[`runner:fileio`]=[]string{`executable`,`fin`,`config`}
	ouLabel[`runner:fileio`]=[]string{`fout`,`stderr`,`judgerlog`}
	inLabel[`runner:stdio`]=[]string{`executable`,`stdin`,`limit`}
	ouLabel[`runner:stdio`]=[]string{`stdout`,`stderr`,`judgerlog`}
}
