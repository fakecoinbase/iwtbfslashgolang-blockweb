// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/node/version.proto

package node

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Version struct {
	Version              int32    `protobuf:"varint,1,opt,name=version,proto3" json:"version,omitempty"`
	BestHeight           uint64   `protobuf:"varint,2,opt,name=bestHeight,proto3" json:"bestHeight,omitempty"`
	AddressFrom          string   `protobuf:"bytes,3,opt,name=addressFrom,proto3" json:"addressFrom,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Version) Reset()         { *m = Version{} }
func (m *Version) String() string { return proto.CompactTextString(m) }
func (*Version) ProtoMessage()    {}
func (*Version) Descriptor() ([]byte, []int) {
	return fileDescriptor_1aa80882f47c9d16, []int{0}
}

func (m *Version) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Version.Unmarshal(m, b)
}
func (m *Version) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Version.Marshal(b, m, deterministic)
}
func (m *Version) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Version.Merge(m, src)
}
func (m *Version) XXX_Size() int {
	return xxx_messageInfo_Version.Size(m)
}
func (m *Version) XXX_DiscardUnknown() {
	xxx_messageInfo_Version.DiscardUnknown(m)
}

var xxx_messageInfo_Version proto.InternalMessageInfo

func (m *Version) GetVersion() int32 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *Version) GetBestHeight() uint64 {
	if m != nil {
		return m.BestHeight
	}
	return 0
}

func (m *Version) GetAddressFrom() string {
	if m != nil {
		return m.AddressFrom
	}
	return ""
}

func init() {
	proto.RegisterType((*Version)(nil), "farmer.Version")
}

func init() { proto.RegisterFile("proto/node/version.proto", fileDescriptor_1aa80882f47c9d16) }

var fileDescriptor_1aa80882f47c9d16 = []byte{
	// 145 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x28, 0x28, 0xca, 0x2f,
	0xc9, 0xd7, 0xcf, 0xcb, 0x4f, 0x49, 0xd5, 0x2f, 0x4b, 0x2d, 0x2a, 0xce, 0xcc, 0xcf, 0xd3, 0x03,
	0x0b, 0x09, 0x71, 0xa4, 0xe4, 0x15, 0xc7, 0x17, 0xa7, 0xa6, 0xa6, 0x28, 0xa5, 0x72, 0xb1, 0x87,
	0x41, 0xa4, 0x84, 0x24, 0xb8, 0xd8, 0xa1, 0xaa, 0x24, 0x18, 0x15, 0x18, 0x35, 0x58, 0x83, 0x60,
	0x5c, 0x21, 0x39, 0x2e, 0xae, 0xa4, 0xd4, 0xe2, 0x12, 0x8f, 0xd4, 0xcc, 0xf4, 0x8c, 0x12, 0x09,
	0x26, 0x05, 0x46, 0x0d, 0x96, 0x20, 0x24, 0x11, 0x21, 0x05, 0x2e, 0xee, 0xc4, 0x94, 0x94, 0xa2,
	0xd4, 0xe2, 0x62, 0xb7, 0xa2, 0xfc, 0x5c, 0x09, 0x66, 0x05, 0x46, 0x0d, 0xce, 0x20, 0x64, 0x21,
	0x27, 0xfe, 0x28, 0xde, 0xcc, 0xbc, 0x92, 0xd4, 0xa2, 0xbc, 0xc4, 0x1c, 0xb0, 0x7b, 0x92, 0xd8,
	0xc0, 0x0e, 0x31, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0x6c, 0x84, 0x74, 0xb6, 0xa4, 0x00, 0x00,
	0x00,
}
