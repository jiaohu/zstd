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
}

// FrameHeader detail of frame header content
//  --------------------------------------
//  | Frame_Header_Descriptor | 1 byte    |
//  | [Window_Descriptor]     | 0-1 byte  |
//  | [Dictionary_ID]         | 0-4 bytes |
//  | [Frame_Content_Size]    | 0-8 bytes |
//  ---------------------------------------
// Detail of
// Frame_Header_Descriptor:
//    It describes which other fields are present. Decoding this byte is enough to tell the size of Frame_Header.
//    Frame_Header_Descriptor contain 8 bits, here is the mean of each bit, bit 7 is the highest bit, 0 is the lowest one.
//    ---------------------------------------
//    | bit number | filed_name              |
//    ----------------------------------------
//    | 7-6        | Frame_Content_Size_Flag |
//    | 5          | Single_Segment_Flag     |
//    | 4          | (unused)                |
//    | 3          | (reserved)              |
//    | 2          | Content_Checksum_Flag   |
//    | 1-0        | Dictionary_ID_Flag      |
//    ----------------------------------------
//    So Frame_Content_Size_Flag will be 0, 1, 2, 3, and it decides FCS_Field_Size bytes (2^Frame_Content_Size_Flag).
//    Must care for Frame_Content_Size_Flag == 0
//      1). FCS_Field_Size depends on Single_Segment_Flag, if  Single_Segment_Flag
//    is set, FCS_Field_Size is 1, otherwise 0;
//      2). Frame_Content_Size is not provided.
//
//    If Single_Segment_Flag is set, data must be regenerated within a single continuous memory segment.
//
//
//
type FrameHeader struct {
	FrameHeaderDescriptor byte  // 1 byte
	WindowDescriptor      byte  // 0-1 byte
	DictionaryId          int32 // 0-4 bytes
	FrameContentSize      int64 // 0-8 bytes
}

type DataBlock struct {
}

func Encode() {

}

func Decode() {

}
