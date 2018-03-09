// Code generated by protoc-gen-go. DO NOT EDIT.
// source: sms.proto

/*
Package smsAli is a generated protocol buffer package.

It is generated from these files:
	sms.proto

It has these top-level messages:
	SendSms_Request
	SendSms_Response
*/
package smsAli

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

type SendSms_Request struct {
	VerifyCode  string `protobuf:"bytes,1,opt,name=verifyCode" json:"verifyCode,omitempty"`
	PhoneNumber string `protobuf:"bytes,2,opt,name=phoneNumber" json:"phoneNumber,omitempty"`
	UserId      string `protobuf:"bytes,3,opt,name=userId" json:"userId,omitempty"`
}

func (m *SendSms_Request) Reset()                    { *m = SendSms_Request{} }
func (m *SendSms_Request) String() string            { return proto.CompactTextString(m) }
func (*SendSms_Request) ProtoMessage()               {}
func (*SendSms_Request) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *SendSms_Request) GetVerifyCode() string {
	if m != nil {
		return m.VerifyCode
	}
	return ""
}

func (m *SendSms_Request) GetPhoneNumber() string {
	if m != nil {
		return m.PhoneNumber
	}
	return ""
}

func (m *SendSms_Request) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

type SendSms_Response struct {
}

func (m *SendSms_Response) Reset()                    { *m = SendSms_Response{} }
func (m *SendSms_Response) String() string            { return proto.CompactTextString(m) }
func (*SendSms_Response) ProtoMessage()               {}
func (*SendSms_Response) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func init() {
	proto.RegisterType((*SendSms_Request)(nil), "Dolphin.Protocol.SendSms_Request")
	proto.RegisterType((*SendSms_Response)(nil), "Dolphin.Protocol.SendSms_Response")
}

func init() { proto.RegisterFile("sms.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 167 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2c, 0xce, 0x2d, 0xd6,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x12, 0x70, 0xc9, 0xcf, 0x29, 0xc8, 0xc8, 0xcc, 0xd3, 0x0b,
	0x00, 0x71, 0x93, 0xf3, 0x73, 0x94, 0xb2, 0xb9, 0xf8, 0x83, 0x53, 0xf3, 0x52, 0x82, 0x73, 0x8b,
	0xe3, 0x83, 0x52, 0x0b, 0x4b, 0x53, 0x8b, 0x4b, 0x84, 0xe4, 0xb8, 0xb8, 0xca, 0x52, 0x8b, 0x32,
	0xd3, 0x2a, 0x9d, 0xf3, 0x53, 0x52, 0x25, 0x18, 0x15, 0x18, 0x35, 0x38, 0x83, 0x90, 0x44, 0x84,
	0x14, 0xb8, 0xb8, 0x0b, 0x32, 0xf2, 0xf3, 0x52, 0xfd, 0x4a, 0x73, 0x93, 0x52, 0x8b, 0x24, 0x98,
	0xc0, 0x0a, 0x90, 0x85, 0x84, 0xc4, 0xb8, 0xd8, 0x4a, 0x8b, 0x53, 0x8b, 0x3c, 0x53, 0x24, 0x98,
	0xc1, 0x92, 0x50, 0x9e, 0x92, 0x10, 0x97, 0x00, 0xc2, 0xb2, 0xe2, 0x82, 0xfc, 0xbc, 0xe2, 0x54,
	0x27, 0xa1, 0x28, 0x01, 0x3d, 0x3d, 0xfd, 0xdc, 0xfc, 0x94, 0xd2, 0x9c, 0x54, 0xfd, 0xe2, 0xdc,
	0x62, 0xc7, 0x9c, 0xcc, 0x24, 0x36, 0xb0, 0x6b, 0x8d, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0xbb,
	0xf6, 0x01, 0xfa, 0xba, 0x00, 0x00, 0x00,
}
