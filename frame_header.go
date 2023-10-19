package zstd

import "errors"

// FrameHeader detail of frame header content
//
//	--------------------------------------
//	| Frame_Header_Descriptor | 1 byte    |
//	| [Window_Descriptor]     | 0-1 byte  |
//	| [Dictionary_ID]         | 0-4 bytes |
//	| [Frame_Content_Size]    | 0-8 bytes |
//	---------------------------------------
//
// Detail of
// Frame_Header_Descriptor:
//
//	It describes which other fields are present. Decoding this byte is enough to tell the size of Frame_Header.
//	Frame_Header_Descriptor contain 8 bits, here is the mean of each bit, bit 7 is the highest bit, 0 is the lowest one.
//	---------------------------------------
//	| bit number | filed_name              |
//	----------------------------------------
//	| 7-6        | Frame_Content_Size_Flag |
//	| 5          | Single_Segment_Flag     |
//	| 4          | (unused)                |
//	| 3          | (reserved)              |
//	| 2          | Content_Checksum_Flag   |
//	| 1-0        | Dictionary_ID_Flag      |
//	----------------------------------------
//	So Frame_Content_Size_Flag will be 0, 1, 2, 3, and it decides FCS_Field_Size bytes (2^Frame_Content_Size_Flag).
//	Must care for Frame_Content_Size_Flag == 0
//	  1). FCS_Field_Size depends on Single_Segment_Flag, if  Single_Segment_Flag
//	is set, FCS_Field_Size is 1, otherwise 0;
//	  2). Frame_Content_Size is not provided.
//
//	If Single_Segment_Flag is set, data must be regenerated within a single continuous memory segment.
type FrameHeader struct {
	FrameHeaderDescriptor byte   // 1 byte
	WindowDescriptor      byte   // 0-1 byte
	DictionaryId          uint32 // 0-4 bytes
	FrameContentSize      uint64 // 0-8 bytes
}

func (f *FrameHeader) FrameContentSizeFlag() byte {
	return f.FrameHeaderDescriptor >> 6
}

func (f *FrameHeader) SingleSegmentFlag() byte {
	return (f.FrameHeaderDescriptor & 0b00100000) >> 5
}

// not set
func (f *FrameHeader) Unused() byte {
	return (f.FrameHeaderDescriptor & 0b00010000) >> 4
}

// not set
func (f *FrameHeader) Reserved() byte {
	return (f.FrameHeaderDescriptor & 0b00001000) >> 3
}

// If this flag is set, a 32-bit Content_Checksum will be present at the frame's end
func (f *FrameHeader) ContentChecksumFlag() byte {
	return (f.FrameHeaderDescriptor & 0b0000100) >> 2
}

func (f *FrameHeader) DictionaryIDFlag() byte {
	return f.FrameHeaderDescriptor & 0b0000011
}

// FCSFieldSize the decompressed data size, which decide the range of FrameContentSize
//
//	--------------------------------------
//	| fcs filed size    |   Range        |
//	--------------------------------------
//	|  0                |  unknown       |
//	|  1                | 0 - 255        |
//	|  2                | 256 - 65791    |
//	|  4                | 0 - 2^32 - 1   |
//	|  8                | 0 - 2^64 - 1   |
//	--------------------------------------
//
// When FCS_Field_Size is 1, 4, or 8 bytes, the value is read directly.
// When FCS_Field_Size is 2, the offset of 256 is added.
// It's allowed to represent a small size (for example, 18) using any compatible variant
func (f *FrameHeader) FCSFieldSize() (int, error) {
	res := f.FrameContentSizeFlag()
	switch res {
	case 0:
		if f.SingleSegmentFlag() == byte(1) {
			return 1, nil
		} else {
			return 0, nil
		}
	case 1:
		return 2, nil
	case 2:
		return 4, nil
	case 4:
		return 8, nil
	default:
		return 0, errors.New("error frame content size flag")
	}
}

func (f *FrameHeader) DIDFieldSize() (int, error) {
	res := f.DictionaryIDFlag()
	switch res {
	case 0:
		return 0, nil
	case 1:
		return 1, nil
	case 2:
		return 2, nil
	case 3:
		return 4, nil
	default:
		return 0, errors.New("error dictionary id flag")
	}
}

// WindowSize the minimum memory buffer size required to decompress a frame
// if Single_Segment_Flag is set, Window_Descriptor is not present, Window_Size is
// Frame_Content_Size, which can be any value from 0 to 2^64-1 bytes
func (f *FrameHeader) WindowSize() uint64 {
	if f.SingleSegmentFlag() == byte(1) {
		return f.FrameContentSize
	}
	exponent := f.WindowDescriptor >> 3
	mantissa := f.WindowDescriptor & 0b111

	windowLog := uint64(10 + exponent)
	windowBase := uint64(1 << windowLog)
	windowAdd := (windowLog / 8) * uint64(mantissa)
	// The minimum Window_Size is 1 KB. The maximum Window_Size is (1<<41) + 7*(1<<38) bytes, which is 3.75 TB.
	// In general, larger Window_Size values tend to improve the compression ratio, but at the cost of increased memory usage.
	// To properly decode compressed data, a decoder will need to allocate a buffer of at least Window_Size bytes.
	// To protect decoders, a decoder is allowed to reject a frame that requests a memory size beyond the size.

	// TODO it might be more excellent when the size to 8MB, check it.
	return windowBase + windowAdd
}

func (f *FrameHeader) DictionaryID() (uint32, error) {
	// One byte can represent an ID 0-255; 2 bytes can represent an ID 0-65535; 4 bytes can represent an ID 0-4294967295. Format is little-endian.
	if f.DictionaryId <= 32767 || f.DictionaryId >= (1<<31) {
		return 0, errors.New("out of dictionary id")
	}

	return f.DictionaryId, nil
}
