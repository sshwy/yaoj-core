package processors

import (
	"fmt"
	"time"

	"github.com/sshwy/yaoj-core/pkg/private/judger"
)

// `s` contains a series of number seperated by space, denoting
// real time (ms), cpu time (ms), virtual memory (byte), real memory (byte),
// stack memory (byte), output limit (byte), fileno limitation respectively.
func parseJudgerLimit(s string) ([]judger.OptionProvider, error) {
	var rt, ct, vm, rm, sm, ol, fl int
	if _, err := fmt.Sscanf(s, "%d%d%d%d%d%d%d", &rt, &ct, &vm, &rm, &sm, &ol, &fl); err != nil {
		return nil, err
	}
	options := []judger.OptionProvider{}
	if rt > 0 {
		options = append(options, judger.WithRealTime(time.Millisecond*time.Duration(rt)))
	}
	if ct > 0 {
		options = append(options, judger.WithCpuTime(time.Millisecond*time.Duration(ct)))
	}
	if vm > 0 {
		options = append(options, judger.WithVirMemory(judger.ByteValue(vm)))
	}
	if rm > 0 {
		options = append(options, judger.WithRealMemory(judger.ByteValue(rm)))
	}
	if sm > 0 {
		options = append(options, judger.WithStack(judger.ByteValue(sm)))
	}
	if ol > 0 {
		options = append(options, judger.WithOutput(judger.ByteValue(ol)))
	}
	if fl > 0 {
		options = append(options, judger.WithFileno(fl))
	}
	return options, nil
}
