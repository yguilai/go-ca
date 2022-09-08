package constantpool

import "github.com/yguilai/go-ca/reader"

type ConstantNameAndTypeInfo struct {
	nameIndex       uint16
	descriptorIndex uint16
}

func NewConstantNameAndTypeInfo() *ConstantNameAndTypeInfo {
	return &ConstantNameAndTypeInfo{}
}

func (self *ConstantNameAndTypeInfo) readInfo(reader *reader.ClassReader) {
	self.nameIndex = reader.ReadUInt16()
	self.descriptorIndex = reader.ReadUInt16()
}

func (self *ConstantNameAndTypeInfo) String() string {
	return ""
}
