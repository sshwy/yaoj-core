package processor

var processors map[string]Processor = make(map[string]Processor)

func Get(name string) Processor {
	return processors[name]
}

func GetAll() map[string]Processor {
	return processors
}

// register a processor to system
func Register(name string, proc Processor) {
	processors[name] = proc
}

func init() {
	Register("checker:hcmp", CheckerHcmp{})
	Register("checker:testlib", CheckerTestlib{})
	Register("compiler", Compiler{})
	Register("generator", Generator{})
	Register("generator:testlib", GeneratorTestlib{})
	Register("runner:fileio", RunnerFileio{})
	Register("runner:stdio", RunnerStdio{})
}
