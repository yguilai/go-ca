package reader

import "encoding/binary"

type ClassReader struct {
	data []byte
}

func NewClassReader(data []byte) *ClassReader {
	return &ClassReader{data}
}

func (r *ClassReader) ReadUInt8() uint8 {
	val := r.data[0]
	r.data = r.data[1:]
	return val
}

func (r *ClassReader) ReadUInt16() uint16 {
	val := binary.BigEndian.Uint16(r.data)
	r.data = r.data[2:]
	return val
}

func (r *ClassReader) ReadUInt32() uint32 {
	val := binary.BigEndian.Uint32(r.data)
	r.data = r.data[4:]
	return val
}

func (r *ClassReader) ReadUInt64() uint64 {
	val := binary.BigEndian.Uint64(r.data)
	r.data = r.data[8:]
	return val
}

func (r *ClassReader) ReadBytes(n uint32) []byte {
	val := r.data[:n]
	r.data = r.data[n:]
	return val
}

func (r *ClassReader) Size() int {
	return len(r.data)
}
