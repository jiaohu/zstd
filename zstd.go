package zstd

// reference from https://datatracker.ietf.org/doc/html/rfc8878#name-compression-algorithm

const FrameMagic = 0xFD2FB528

const ContentChecksumFlag = false

const BlockMaximumSize = 128 * 1024

// ZFrame is a single struct
//
//	----------------------------------
//	| Magic_Number       |  4 bytes  |
//	| Frame_Header       | 2-14 bytes|
//	| Data_Block         | n bytes   |
//	| [More Data_Blocks] |           |
//	| [Content_Checksum] | 4 bytes   |
//	----------------------------------
type ZFrame struct {
	magicNumber [4]byte
	frameHeader FrameHeader
	block       DataBlock
	checksum    [4]byte
}

func (z *ZFrame) IsChecksumPresent() bool {
	return z.frameHeader.ContentChecksumFlag() == byte(1)
}

func (z *ZFrame) GetBlockMaximumSize() int {
	windowSize := z.frameHeader.WindowSize()
	if windowSize > BlockMaximumSize {
		return BlockMaximumSize
	}
	return int(windowSize)
}

func (z *ZFrame) checkBlocks() bool {
	return z.block.BlockSize() < uint32(z.GetBlockMaximumSize())
}
