package processors

import (
	"log"
	"os"

	"github.com/sshwy/yaoj-core/pkg/processor"
)

type Processor = processor.Processor

type Result = processor.Result

var logger = log.New(os.Stderr, "[processors] ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix)
