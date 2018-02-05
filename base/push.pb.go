// Code generated by protoc-gen-go. DO NOT EDIT.
// source: push.proto

/*
Package base is a generated protocol buffer package.

It is generated from these files:
	push.proto

It has these top-level messages:
	PushContent
*/
package base

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type PushContent struct {
	// 推送内容
	Content []byte `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
}

func (m *PushContent) Reset()                    { *m = PushContent{} }
func (m *PushContent) String() string            { return proto.CompactTextString(m) }
func (*PushContent) ProtoMessage()               {}
func (*PushContent) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *PushContent) GetContent() []byte {
	if m != nil {
		return m.Content
	}
	return nil
}

func init() {
	proto.RegisterType((*PushContent)(nil), "base.PushContent")
}

func init() { proto.RegisterFile("push.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 78 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2a, 0x28, 0x2d, 0xce,
	0xd0, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x49, 0x4a, 0x2c, 0x4e, 0x55, 0x52, 0xe7, 0xe2,
	0x0e, 0x28, 0x2d, 0xce, 0x70, 0xce, 0xcf, 0x2b, 0x49, 0xcd, 0x2b, 0x11, 0x92, 0xe0, 0x62, 0x4f,
	0x86, 0x30, 0x25, 0x18, 0x15, 0x18, 0x35, 0x78, 0x82, 0x60, 0xdc, 0x24, 0x36, 0xb0, 0x2e, 0x63,
	0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x78, 0x03, 0xb8, 0x1b, 0x43, 0x00, 0x00, 0x00,
}