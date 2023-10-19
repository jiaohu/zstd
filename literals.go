package zstd

import "errors"

type LiteralBlockType = byte

const (
	// Literals are stored uncompressed. Literals_Section_Content is Regenerated_Size
	RawLiteralBlock LiteralBlockType = iota
	// Literals consist of a single-byte value repeated Regenerated_Size times. Literals_Section_Content is 1.
	RLELiteralBlock
	// This is a standard Huffman-compressed block, starting with a Huffman tree description. See details below. Literals_Section_Content is Compressed_Size.
	CompressedLiteralBlock
	// This is a Huffman-compressed block, using the Huffman tree from the previous Compressed_Literals_Block or a dictionary if there is no previous Huffman-compressed literals block.
	// Huffman_Tree_Description will be skipped.
	// Literals_Section_Content is Compressed_Size.
	TreelessLiteralBlock
)

type Literal struct {
	// header include:
	//  	Literals_block_type 2 bits
	//		Size format 1-2 bits
	// 		Regenerated_Size 5-20 bits
	// 		[Compressed_Size]  0-18 bits
	header                   []byte // range 1 to 5 bytes
	headerSize               int    // header total size for byte
	headerRegSize            int
	headerRegSizeBits        int // header size format bits len
	headerCompressedSizeBits int // header compressed size bits len
	// only present when type is CompressedLiteralBlock
	// it determine where streams begin Total_Streams_Size = Compressed_Size - Huffman_Tree_Description_Size
	huffmanTreeDescription interface{}
	// The Jump_Table is only present when there are 4 Huffman-coded streams(only streamSize = 4)
	// jumpTable is 6 bytes long and consists of three 2-byte, describe the size of first 3 streams.
	// stream4 size is computed
	// Stream4_Size = Total_Streams_Size - 6 - Stream1_Size - Stream2_Size - Stream3_Size
	jumpTable  [6]byte
	streams    [][]byte
	streamSize int // streams size, stream1 or [stream1, ..., stream4]
}

func (l *Literal) GetLiteralBlockType() LiteralBlockType {
	return l.header[0] & 0b11
}

func (l *Literal) GetSizeFormat() byte {
	return (l.header[0] >> 2) & 0x3
}

// GetRegeneratedSize
func (l *Literal) InitHeaderSection() error {
	switch l.GetLiteralBlockType() {
	// it's only necessary to decode Regenerated_Size. There is no Compressed_Size field.
	case RawLiteralBlock, RLELiteralBlock:
		l.streamSize = 1
		switch l.GetSizeFormat() {
		case 0b00, 0b10:
			// Regenerated_Size use 5 bits (value 0-31)
			l.headerSize = 1
			l.headerRegSizeBits = 5
			l.headerRegSize = int(l.header[0] >> 3)
			return nil
		case 0b01:
			// Regenerated_Size use 12 bits (value 0-4095)
			l.headerSize = 2
			l.headerRegSizeBits = 12
			l.headerRegSize = int(l.header[0]>>4) + int(l.header[1]<<4)
			return nil
		case 0b11:
			// Regenerated_Size uses 20 bits (values 0-1048575)
			l.headerSize = 3
			l.headerRegSizeBits = 20
			l.headerRegSize = int(l.header[0]>>4) + int(l.header[1]<<4) + int(l.header[2])<<12
			return nil
		default:
			return errors.New("error size format")
		}
	// it's required to decode both Compressed_Size and Regenerated_Size (the decompressed size). It's also necessary to decode the number of streams (1 or 4).
	case CompressedLiteralBlock, TreelessLiteralBlock:
		switch l.GetSizeFormat() {
		case 0b00:
			l.streamSize = 1
			// (values 0-1023)
			l.headerRegSizeBits = 10
			l.headerCompressedSizeBits = 10
			l.headerSize = 3
			return nil
		case 0b01:
			l.streamSize = 4
			// (values 0-1023)
			l.headerRegSizeBits = 10
			l.headerCompressedSizeBits = 10
			l.headerSize = 3
			return nil
		case 0b10:
			l.streamSize = 4
			// (values 0-16383)
			l.headerRegSizeBits = 14
			l.headerCompressedSizeBits = 14
			l.headerSize = 4
			return nil
		case 0b11:
			l.streamSize = 4
			// (values 0-262143)
			l.headerRegSizeBits = 18
			l.headerCompressedSizeBits = 18
			l.headerSize = 5
			return nil
		default:
			return errors.New("error size format")
		}
	default:
		return errors.New("error literal block type")
	}
}
