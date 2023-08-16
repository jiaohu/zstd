package zstd

type SkipFrame struct {
	// Value: 0x184D2A5?, which means any value from 0x184D2A50 to 0x184D2A5F.
	// All 16 values are valid to identify a skippable frame.
	magicNumber [4]byte
	frameSize   uint32
	userData    []byte
}
