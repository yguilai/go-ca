package classfile

import (
	"fmt"
	"github.com/yguilai/go-ca/accessflags"
	"github.com/yguilai/go-ca/attribute"
	"github.com/yguilai/go-ca/constantpool"
	"github.com/yguilai/go-ca/reader"
	"strconv"
	"strings"
)

type ClassFile struct {
	magic             uint32
	minorVersion      uint16
	majorVersion      uint16
	constantPoolCount uint16
	constantPool      []constantpool.ConstantInfo
	accessFlags       uint16
	thisClass         uint16
	superClass        uint16
	interfacesCount   uint16
	interfaces        []uint16
	fieldsCount       uint16
	fields            []FieldInfo
	methodsCount      uint16
	methods           []MethodInfo
	attributesCount   uint16
	attributes        []attribute.AttributeInfo
}

type FieldInfo struct {
	accessFlags     uint16
	nameIndex       uint16
	descriptorIndex uint16
	attributesCount uint16
	attributes      []attribute.AttributeInfo
}

type MethodInfo struct {
	accessFlags     uint16
	nameIndex       uint16
	descriptorIndex uint16
	attributesCount uint16
	attributes      []attribute.AttributeInfo
}

func NewClassFile() *ClassFile {
	return &ClassFile{}
}

func (cf *ClassFile) Parse(reader *reader.ClassReader) (*ClassFile, error) {
	cf.ReadMagic(reader)
	cf.ReadMinorVersion(reader)
	cf.ReadMajorVersion(reader)

	// parse constant pool info
	cf.ReadConstantPoolCount(reader)
	cf.constantPool = constantpool.Parse(reader, cf.constantPoolCount)

	cf.ReadAccessFlags(reader)
	cf.ReadThisClass(reader)
	cf.ReadSuperClass(reader)

	cf.ReadInterfacesCount(reader)
	cf.ReadInterfaces(reader, cf.interfacesCount)

	// parse fields info
	cf.ReadFieldsCount(reader)
	cf.ParseFields(reader, cf.fieldsCount)

	// parse methods info
	cf.ReadMethodsCount(reader)
	cf.ParseMethods(reader, cf.methodsCount)

	// parse attributes info
	cf.ReadAttributesCount(reader)
	cf.ParseClassAttributes(reader, cf.attributesCount)

	return cf, nil
}

func (cf *ClassFile) ParseFields(reader *reader.ClassReader, count uint16) (*ClassFile, error) {
	cf.fields = make([]FieldInfo, 0)
	for index := 0; index < int(count); index++ {
		field := FieldInfo{}
		field.accessFlags = reader.ReadUInt16()
		field.nameIndex = reader.ReadUInt16()
		field.descriptorIndex = reader.ReadUInt16()
		field.attributesCount = reader.ReadUInt16()
		field.attributes = attribute.ReadAttributes(reader, cf.constantPool, field.attributesCount)
		cf.fields = append(cf.fields, field)
	}
	return cf, nil
}

func (cf *ClassFile) ParseMethods(reader *reader.ClassReader, count uint16) (*ClassFile, error) {
	cf.methods = make([]MethodInfo, 0)
	for index := 0; index < int(count); index++ {
		method := MethodInfo{}
		method.accessFlags = reader.ReadUInt16()
		method.nameIndex = reader.ReadUInt16()
		method.descriptorIndex = reader.ReadUInt16()
		method.attributesCount = reader.ReadUInt16()
		method.attributes = attribute.ReadAttributes(reader, cf.constantPool, method.attributesCount)
		cf.methods = append(cf.methods, method)
	}
	return cf, nil
}

// Print class structure
func (cf *ClassFile) ViewClass() {
	fmt.Printf("magic: %x\n", cf.magic)
	fmt.Printf("minor version: %d\n", cf.minorVersion)
	fmt.Printf("major version: %d\n", cf.majorVersion)
	fmt.Printf("constant pool count: %d\n", cf.constantPoolCount)
	fmt.Printf("access flags: %s\n", cf.AccessFlags(cf.accessFlags))
	fmt.Printf("this class: %s\n", cf.ThisClass())
	fmt.Printf("super class: %s\n", cf.SuperClass())
	fmt.Printf("interfaces: %s\n", cf.Interfaces())
	fmt.Printf("fields count: %d\n", cf.fieldsCount)
	cf.ViewFields()
	fmt.Printf("method count: %d\n", cf.methodsCount)
	cf.ViewMethods()
}

// Print class fields info
func (cf *ClassFile) ViewFields() {
	fields := cf.fields
	for _, field := range fields {
		fmt.Println()
		fmt.Printf("field name: %s\n", cf.GetUtf8Value(field.nameIndex))
		fmt.Printf("access flags: %s\n", accessflags.GetFieldAccessFlags(field.accessFlags))
		fmt.Printf("descriptor index: %s\n", cf.GetUtf8Value(field.descriptorIndex))
	}
	fmt.Println()
}

// Print class methods info
func (cf *ClassFile) ViewMethods() {
	methods := cf.methods
	for _, method := range methods {
		fmt.Println()
		fmt.Printf("method name: %s\n", cf.GetUtf8Value(method.nameIndex))
		fmt.Printf("access flags: %s\n", accessflags.GetMethodAccessFlags(method.accessFlags))
		fmt.Printf("descriptor index: %s\n", cf.GetUtf8Value(method.descriptorIndex))
	}
}

// Get constantpool's CONSTANT_Utf8_info value.
func (cf *ClassFile) GetUtf8Value(index uint16) string {
	if index <= 0 || index >= cf.constantPoolCount {
		panic("GetUtf8Value: Index Out of Bounds")
	}
	val := cf.constantPool[index-1].String()
	return val
}

// Parse class level attributes.
func (cf *ClassFile) ParseClassAttributes(reader *reader.ClassReader, count uint16) (*ClassFile, error) {
	cf.attributes = attribute.ReadAttributes(reader, cf.constantPool, count)
	return cf, nil
}

func (cf *ClassFile) ReadMagic(reader *reader.ClassReader) (*ClassFile, error) {
	val := reader.ReadUInt32()
	if val != 0xCAFEBABE {
		panic("java.lang.ClassFormatCheckError: magic!")
	}
	cf.magic = val
	return cf, nil
}

func (cf *ClassFile) ReadMinorVersion(reader *reader.ClassReader) (*ClassFile, error) {
	val := reader.ReadUInt16()
	cf.majorVersion = val
	return cf, nil
}

func (cf *ClassFile) ReadMajorVersion(reader *reader.ClassReader) (*ClassFile, error) {
	val := reader.ReadUInt16()
	cf.majorVersion = val
	return cf, nil
}

func (cf *ClassFile) ReadConstantPoolCount(reader *reader.ClassReader) (*ClassFile, error) {
	val := reader.ReadUInt16()
	cf.constantPoolCount = val
	return cf, nil
}

func (cf *ClassFile) ReadAccessFlags(reader *reader.ClassReader) (*ClassFile, error) {
	val := reader.ReadUInt16()
	cf.accessFlags = val
	return cf, nil
}

func (cf *ClassFile) ReadThisClass(reader *reader.ClassReader) (*ClassFile, error) {
	val := reader.ReadUInt16()
	cf.thisClass = val
	return cf, nil
}

func (cf *ClassFile) ReadSuperClass(reader *reader.ClassReader) (*ClassFile, error) {
	val := reader.ReadUInt16()
	cf.superClass = val
	return cf, nil
}

func (cf *ClassFile) ReadInterfacesCount(reader *reader.ClassReader) (*ClassFile, error) {
	val := reader.ReadUInt16()
	cf.interfacesCount = val
	return cf, nil
}

func (cf *ClassFile) ReadInterfaces(reader *reader.ClassReader, count uint16) (*ClassFile, error) {
	var index uint16
	cf.interfaces = make([]uint16, 0)
	for index = 0; index < count; index++ {
		cf.interfaces = append(cf.interfaces, reader.ReadUInt16())
	}
	return cf, nil
}

func (cf *ClassFile) ReadFieldsCount(reader *reader.ClassReader) (*ClassFile, error) {
	val := reader.ReadUInt16()
	cf.fieldsCount = val
	return cf, nil
}

func (cf *ClassFile) ReadMethodsCount(reader *reader.ClassReader) (*ClassFile, error) {
	val := reader.ReadUInt16()
	cf.methodsCount = val
	return cf, nil
}

func (cf *ClassFile) ReadAttributesCount(reader *reader.ClassReader) (*ClassFile, error) {
	val := reader.ReadUInt16()
	cf.attributesCount = val
	return cf, nil
}

func (cf *ClassFile) Magic() uint32 {
	return cf.magic
}

func (cf *ClassFile) AccessFlags(flags uint16) string {
	return accessflags.GetClassAccessFlags(flags)
}

func (cf *ClassFile) ThisClass() string {
	index, _ := strconv.Atoi(cf.GetUtf8Value(cf.thisClass))
	return cf.GetUtf8Value(uint16(index))
}

func (cf *ClassFile) SuperClass() string {
	index, _ := strconv.Atoi(cf.GetUtf8Value(cf.superClass))
	return cf.GetUtf8Value(uint16(index))
}

func (cf *ClassFile) Interfaces() string {
	var index int
	interfaces := make([]string, 0)
	for index = 0; index < len(cf.interfaces); index++ {
		interIndex, _ := strconv.Atoi(cf.GetUtf8Value(cf.interfaces[index]))
		interfaces = append(interfaces, cf.GetUtf8Value(uint16(interIndex)))
	}
	return strings.Join(interfaces, ", ")
}

func (cf *ClassFile) MethodName(nameIndex uint16) string {
	return cf.GetUtf8Value(nameIndex)
}

// Get class's main method.
func (cf *ClassFile) MainMethod() *MethodInfo {
	methods := cf.methods
	for _, method := range methods {
		methodName := cf.GetUtf8Value(method.nameIndex)
		descriptor := cf.GetUtf8Value(method.descriptorIndex)
		if methodName == "main" && descriptor == "([Ljava/lang/String;)V" {
			return &method
		}
	}
	return nil
}

func (cf *MethodInfo) Attributes() []attribute.AttributeInfo {
	return cf.attributes
}
