package zstd

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
	header                 []byte
	huffmanTreeDescription interface{}
	jumpTable              interface{}
	stream                 [][]byte
}

func (l *Literal) GetLiteralBlockType() LiteralBlockType {
	return l.header[0] & 0b11
}
