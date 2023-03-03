package zstd

import (
	"encoding/binary"
	"fmt"
	"testing"
)

func TestSetLast(t *testing.T) {
	a := uint32(256)
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, a)
	fmt.Println(b)

	c := make([]byte, 4)
	binary.BigEndian.PutUint32(c, a)
	fmt.Println(c)
}
