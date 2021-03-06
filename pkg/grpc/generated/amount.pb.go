// Code generated by protoc-gen-go. DO NOT EDIT.
// source: amount.proto

package generated

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

type Amount struct {
	AssetId              []byte   `protobuf:"bytes,1,opt,name=asset_id,json=assetId,proto3" json:"asset_id,omitempty"`
	Amount               int64    `protobuf:"varint,2,opt,name=amount,proto3" json:"amount,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Amount) Reset()         { *m = Amount{} }
func (m *Amount) String() string { return proto.CompactTextString(m) }
func (*Amount) ProtoMessage()    {}
func (*Amount) Descriptor() ([]byte, []int) {
	return fileDescriptor_53817c82a66f6a1c, []int{0}
}

func (m *Amount) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Amount.Unmarshal(m, b)
}
func (m *Amount) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Amount.Marshal(b, m, deterministic)
}
func (m *Amount) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Amount.Merge(m, src)
}
func (m *Amount) XXX_Size() int {
	return xxx_messageInfo_Amount.Size(m)
}
func (m *Amount) XXX_DiscardUnknown() {
	xxx_messageInfo_Amount.DiscardUnknown(m)
}

var xxx_messageInfo_Amount proto.InternalMessageInfo

func (m *Amount) GetAssetId() []byte {
	if m != nil {
		return m.AssetId
	}
	return nil
}

func (m *Amount) GetAmount() int64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func init() {
	proto.RegisterType((*Amount)(nil), "waves.Amount")
}

func init() { proto.RegisterFile("amount.proto", fileDescriptor_53817c82a66f6a1c) }

var fileDescriptor_53817c82a66f6a1c = []byte{
	// 136 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x49, 0xcc, 0xcd, 0x2f,
	0xcd, 0x2b, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2d, 0x4f, 0x2c, 0x4b, 0x2d, 0x56,
	0xb2, 0xe6, 0x62, 0x73, 0x04, 0x0b, 0x0b, 0x49, 0x72, 0x71, 0x24, 0x16, 0x17, 0xa7, 0x96, 0xc4,
	0x67, 0xa6, 0x48, 0x30, 0x2a, 0x30, 0x6a, 0xf0, 0x04, 0xb1, 0x83, 0xf9, 0x9e, 0x29, 0x42, 0x62,
	0x5c, 0x6c, 0x10, 0xbd, 0x12, 0x4c, 0x0a, 0x8c, 0x1a, 0xcc, 0x41, 0x50, 0x9e, 0x93, 0x3e, 0x97,
	0x54, 0x72, 0x7e, 0xae, 0x1e, 0xd8, 0xa4, 0x82, 0x9c, 0xc4, 0x92, 0xb4, 0xfc, 0xa2, 0x5c, 0x88,
	0xf1, 0x49, 0xa5, 0x69, 0x51, 0x9c, 0xe9, 0xa9, 0x79, 0xa9, 0x45, 0x89, 0x25, 0xa9, 0x29, 0xab,
	0x98, 0x58, 0xc3, 0x41, 0x6a, 0x92, 0xd8, 0xc0, 0x92, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff,
	0xd9, 0xc7, 0xef, 0xb0, 0x8b, 0x00, 0x00, 0x00,
}
