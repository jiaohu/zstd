package zstd

// reference from https://datatracker.ietf.org/doc/html/rfc8878#name-compression-algorithm

const FrameMagic = 0xFD2FB528

const ContentChecksumFlag = false

// ZFrame is a single struct
//  ----------------------------------
//  | Magic_Number       |  4 bytes  |
//  | Frame_Header       | 2-14 bytes|
//  | Data_Block         | n bytes   |
//  | [More Data_Blocks] |           |
//  | [Content_Checksum] | 4 bytes   |
//  ----------------------------------
type ZFrame struct {
	magicNumber uint32
	frameHeader FrameHeader
	block       DataBlock
	checksum    [4]byte
}
