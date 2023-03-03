package zstd

type DataBlock struct {
	header  uint32 // 3 bytes，
	content []byte
}

func (d *DataBlock) LastBlock() bool {
	return d.header&1 == 1
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
