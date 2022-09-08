package constantpool

import "github.com/yguilai/go-ca/reader"

type ConstantMethodRefInfo struct {
	classIndex       uint16
	nameAndTypeIndex uint16
}

func NewConstantMethodRefInfo() *ConstantMethodRefInfo {
	return &ConstantMethodRefInfo{}
}

func (self *ConstantMethodRefInfo) readInfo(reader *reader.ClassReader) {
	self.classIndex = reader.ReadUInt16()
	self.nameAndTypeIndex = reader.ReadUInt16()
}

func (self *ConstantMethodRefInfo) String() string {
	return ""
}
