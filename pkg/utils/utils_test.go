package utils_test

import (
	"bytes"
	"testing"

	"github.com/sshwy/yaoj-core/pkg/utils"
)

func TestChecksum(t *testing.T) {
	sum := utils.ReaderChecksum(bytes.NewReader([]byte("hello1")))
	t.Log(sum)
}
