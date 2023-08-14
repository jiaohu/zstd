package zstd

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/OneOfOne/xxhash"
	"github.com/stretchr/testify/assert"
)

func TestXXHash(t *testing.T) {
	h := xxhash.New64()
	// [104 101 108 108 111 239 70 219 55 81 216 233 153]
	h.Write([]byte("hello"))
	s := make([]byte, 8)
	binary.LittleEndian.PutUint64(s, h.Sum64())
	fmt.Println(s)
}

func TestHash(t *testing.T) {
	h := xxhash.New64()
	fmt.Println(h.Sum([]byte("hello")))
}

func TestDemo(t *testing.T) {
	a := 9
	fmt.Println(a & 0x3)
}

func TestEncode(t *testing.T) {
	//originData := "6080604052348015600f57600080fd5b506004361060285760003560e01c80634f2be91f14602d575b600080fd5b60336047565b604051603e9190605d565b60405180910390f35b60006002905090565b6057816076565b82525050565b6000602082019050607060008301846050565b92915050565b60008160030b905091905056fea26469706673582212204b196e1e349716f746ca8dfcb47956ea08b0cb19b16a80d0fbce7ff226ba655864736f6c63430008070033"

}

func TestFrameHeader_FrameHeaderDescriptor(t *testing.T) {
	frame := FrameHeader{
		FrameHeaderDescriptor: 255,
	}
	t.Run("ContentSizeFlag", func(t *testing.T) {
		assert.Equal(t, byte(3), frame.FrameContentSizeFlag())
	})
	t.Run("SingleSegmentFlag", func(t *testing.T) {
		assert.Equal(t, byte(1), frame.SingleSegmentFlag())
	})
	t.Run("unused", func(t *testing.T) {
		assert.Equal(t, byte(1), frame.Unused())
	})
	t.Run("reserved", func(t *testing.T) {
		assert.Equal(t, byte(1), frame.Reserved())
	})
	t.Run("contentChecksumFlag", func(t *testing.T) {
		assert.Equal(t, byte(1), frame.ContentChecksumFlag())
	})
	t.Run("DictionaryIDFlag", func(t *testing.T) {
		assert.Equal(t, byte(3), frame.DictionaryIDFlag())
	})
}

func Test_U32And4Bytes(t *testing.T) {
	a := "FD2FB528"
	temp, e := hex.DecodeString(a)
	fmt.Println(e)
	fmt.Println(temp)
	a1 := binary.LittleEndian.Uint32(temp)
	fmt.Println(a1)

	b := uint32(682962941)
	s := make([]byte, 4)
	binary.LittleEndian.PutUint32(s, b)
	println(hex.EncodeToString(s))
}
