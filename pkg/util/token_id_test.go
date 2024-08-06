package util_test

import (
	"fmt"
	"github.com/gophero/goal/ciphers"
	"github.com/gophero/goal/conv"
	"goboot/pkg/util"
	"testing"
)

func TestGenerateTxNum(t *testing.T) {
	println(util.GenerateTxNum(1))
	// println(util.UUID32())
}

func TestToken(t *testing.T) {
	ids := []uint{3105, 4839, 9555, 23095, 37254, 80380, 100309, 1321}
	for _, id := range ids {
		prefix := ciphers.MD5(conv.Int64ToStr(int64(id)))
		fmt.Printf("id: %d, prefix: %s\n", id, prefix)
	}
}
