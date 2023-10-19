package zstd

type CompressionMode = byte

const (
	Predefined CompressionMode = iota
	RLE
	FSE
	Repeat
)

type SequenceHeader struct {
	number [3]byte
	symbol byte
}

func (s *SequenceHeader) GetLiteralLengthMode() byte {
	return s.symbol >> 6
}

func (s *SequenceHeader) GetOffsetMode() byte {
	return (s.symbol >> 4) & 0x3
}

func (s *SequenceHeader) GetMatchLengthMode() byte {
	return (s.symbol >> 2) & 0x3
}

func (s *SequenceHeader) GetReserved() byte {
	return s.symbol & 0x3
}

// Sequence A compressed block is a succession of sequences
// Sequences_Section_Size = Block_Size - Literals_Section_Header - Literals_Section_Content.
//
// ---------------------------
// |Sequences_Section_Header |
// |[Literals_Length_Table]  |
// |[Offset_Table]           |
// |[Match_Length_Table]     |
// |bitStream                |
// ---------------------------
type Sequence struct {
	header   SequenceHeader
	bitSteam []byte
}

func (s *Sequence) GetNumberOfSequences() int {
	if s.header.number[0] == 0 {
		// there are no sequences. The sequence section stops here.
		// Decompressed content is defined entirely as Literals_Section content.
		// The FSE tables used in Repeat_Mode are not updated.
		return 0
	} else if s.header.number[0] < 128 {
		return int(s.header.number[0])
	} else if s.header.number[0] < 255 {
		return int(s.header.number[0]-128)<<8 + int(s.header.number[1])
	} else if s.header.number[0] == 255 {
		return int(s.header.number[1]) + int(s.header.number[2])<<8 + 0x7f00
	}
	return 0
}
