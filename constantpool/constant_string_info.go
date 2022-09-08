package constantpool

import "github.com/yguilai/go-ca/reader"

type ConstantStringInfo struct {
	index uint16
}

func NewConstantStringInfo() *ConstantStringInfo {
	return &ConstantStringInfo{}
}

func (self *ConstantStringInfo) readInfo(reader *reader.ClassReader) {
	self.index = reader.ReadUInt16()
}

func (self *ConstantStringInfo) String() string {
	return ""
}
