package zstd

type DataBlock struct {
	header  uint32 // 3 bytes，
	content []byte
}

type BlockType = byte

const (
	// This is an uncompressed block. Block_Content contains Block_Size bytes
	RawBlock BlockType = iota
	// This is a single byte, repeated Block_Size times.
	// Block_Content consists of a single byte. On the decompression side, this byte must be repeated Block_Size times
	RLEBlock
	// This is a compressed block, Block_Size is the length of Block_Content, namely the compressed data.
	CompressedBlock
	// This is not a block. This value cannot be used with the current specification.
	// If such a value is present, it is considered to be corrupt data, and a compliant decoder must reject it.
	Reserved
)

// LastBlock signals whether this block is the last block
func (d *DataBlock) LastBlock() bool {
	return d.header&1 == 1
}

func (d *DataBlock) BlockType() BlockType {
	return BlockType(d.header & 0b11)
}

// When Block_Type is Compressed_Block or Raw_Block, Block_Size is the size of Block_Content (hence excluding Block_Header).
// When Block_Type is RLE_Block, since Block_Content's size is always 1, Block_Size represents the number of times this byte must be repeated
func (d *DataBlock) BlockSize() uint32 {
	return d.header >> 3
}

func (d *DataBlock) SetLastBlock(b bool) {
	if b {
		d.header = d.header | 1
	} else {
		d.header = d.header & (1<<24 - 2)
	}
}

func (d *DataBlock) ToBytes() ([]byte, error) {
	return nil, nil
}
