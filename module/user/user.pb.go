// Code generated by protoc-gen-go. DO NOT EDIT.
// source: user.proto

/*
Package user is a generated protocol buffer package.

It is generated from these files:
	user.proto

It has these top-level messages:
	UserTokenData
	UserData
	User_Login_Request
	User_Login_Response
	User_Register_Request
	User_Register_Response
	User_GetVerifyCode_Request
	User_GetVerifyCode_Response
	User_GetSelfInfo_Request
	User_GetSelfInfo_Response
	User_LogOut_Request
	User_LogOut_Response
	User_RefreshToken_Request
	User_RefreshToken_Response
*/
package user

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

type UserTokenData struct {
	Token        string `protobuf:"bytes,1,opt,name=token" json:"token,omitempty"`
	RefreshToken string `protobuf:"bytes,2,opt,name=refreshToken" json:"refreshToken,omitempty"`
	ExpireAt     int64  `protobuf:"varint,3,opt,name=expireAt" json:"expireAt,omitempty"`
}

func (m *UserTokenData) Reset()                    { *m = UserTokenData{} }
func (m *UserTokenData) String() string            { return proto.CompactTextString(m) }
func (*UserTokenData) ProtoMessage()               {}
func (*UserTokenData) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *UserTokenData) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *UserTokenData) GetRefreshToken() string {
	if m != nil {
		return m.RefreshToken
	}
	return ""
}

func (m *UserTokenData) GetExpireAt() int64 {
	if m != nil {
		return m.ExpireAt
	}
	return 0
}

type UserData struct {
	UserName string `protobuf:"bytes,1,opt,name=UserName" json:"UserName,omitempty"`
	UserID   string `protobuf:"bytes,2,opt,name=UserID" json:"UserID,omitempty"`
}

func (m *UserData) Reset()                    { *m = UserData{} }
func (m *UserData) String() string            { return proto.CompactTextString(m) }
func (*UserData) ProtoMessage()               {}
func (*UserData) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *UserData) GetUserName() string {
	if m != nil {
		return m.UserName
	}
	return ""
}

func (m *UserData) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

type User_Login_Request struct {
	Username string `protobuf:"bytes,1,opt,name=username" json:"username,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password" json:"password,omitempty"`
}

func (m *User_Login_Request) Reset()                    { *m = User_Login_Request{} }
func (m *User_Login_Request) String() string            { return proto.CompactTextString(m) }
func (*User_Login_Request) ProtoMessage()               {}
func (*User_Login_Request) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *User_Login_Request) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *User_Login_Request) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type User_Login_Response struct {
	UserTokenData *UserTokenData `protobuf:"bytes,1,opt,name=user_token_data,json=userTokenData" json:"user_token_data,omitempty"`
}

func (m *User_Login_Response) Reset()                    { *m = User_Login_Response{} }
func (m *User_Login_Response) String() string            { return proto.CompactTextString(m) }
func (*User_Login_Response) ProtoMessage()               {}
func (*User_Login_Response) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *User_Login_Response) GetUserTokenData() *UserTokenData {
	if m != nil {
		return m.UserTokenData
	}
	return nil
}

type User_Register_Request struct {
	Username   string `protobuf:"bytes,1,opt,name=username" json:"username,omitempty"`
	Password   string `protobuf:"bytes,2,opt,name=password" json:"password,omitempty"`
	VerifyCode string `protobuf:"bytes,3,opt,name=verifyCode" json:"verifyCode,omitempty"`
	OpenId     string `protobuf:"bytes,4,opt,name=openId" json:"openId,omitempty"`
}

func (m *User_Register_Request) Reset()                    { *m = User_Register_Request{} }
func (m *User_Register_Request) String() string            { return proto.CompactTextString(m) }
func (*User_Register_Request) ProtoMessage()               {}
func (*User_Register_Request) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *User_Register_Request) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *User_Register_Request) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func (m *User_Register_Request) GetVerifyCode() string {
	if m != nil {
		return m.VerifyCode
	}
	return ""
}

func (m *User_Register_Request) GetOpenId() string {
	if m != nil {
		return m.OpenId
	}
	return ""
}

type User_Register_Response struct {
	Result string `protobuf:"bytes,1,opt,name=result" json:"result,omitempty"`
}

func (m *User_Register_Response) Reset()                    { *m = User_Register_Response{} }
func (m *User_Register_Response) String() string            { return proto.CompactTextString(m) }
func (*User_Register_Response) ProtoMessage()               {}
func (*User_Register_Response) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *User_Register_Response) GetResult() string {
	if m != nil {
		return m.Result
	}
	return ""
}

type User_GetVerifyCode_Request struct {
	PhoneNumber string `protobuf:"bytes,1,opt,name=phoneNumber" json:"phoneNumber,omitempty"`
}

func (m *User_GetVerifyCode_Request) Reset()                    { *m = User_GetVerifyCode_Request{} }
func (m *User_GetVerifyCode_Request) String() string            { return proto.CompactTextString(m) }
func (*User_GetVerifyCode_Request) ProtoMessage()               {}
func (*User_GetVerifyCode_Request) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *User_GetVerifyCode_Request) GetPhoneNumber() string {
	if m != nil {
		return m.PhoneNumber
	}
	return ""
}

type User_GetVerifyCode_Response struct {
}

func (m *User_GetVerifyCode_Response) Reset()                    { *m = User_GetVerifyCode_Response{} }
func (m *User_GetVerifyCode_Response) String() string            { return proto.CompactTextString(m) }
func (*User_GetVerifyCode_Response) ProtoMessage()               {}
func (*User_GetVerifyCode_Response) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

type User_GetSelfInfo_Request struct {
}

func (m *User_GetSelfInfo_Request) Reset()                    { *m = User_GetSelfInfo_Request{} }
func (m *User_GetSelfInfo_Request) String() string            { return proto.CompactTextString(m) }
func (*User_GetSelfInfo_Request) ProtoMessage()               {}
func (*User_GetSelfInfo_Request) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

type User_GetSelfInfo_Response struct {
	UserData *UserData `protobuf:"bytes,1,opt,name=userData" json:"userData,omitempty"`
}

func (m *User_GetSelfInfo_Response) Reset()                    { *m = User_GetSelfInfo_Response{} }
func (m *User_GetSelfInfo_Response) String() string            { return proto.CompactTextString(m) }
func (*User_GetSelfInfo_Response) ProtoMessage()               {}
func (*User_GetSelfInfo_Response) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *User_GetSelfInfo_Response) GetUserData() *UserData {
	if m != nil {
		return m.UserData
	}
	return nil
}

type User_LogOut_Request struct {
}

func (m *User_LogOut_Request) Reset()                    { *m = User_LogOut_Request{} }
func (m *User_LogOut_Request) String() string            { return proto.CompactTextString(m) }
func (*User_LogOut_Request) ProtoMessage()               {}
func (*User_LogOut_Request) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

type User_LogOut_Response struct {
}

func (m *User_LogOut_Response) Reset()                    { *m = User_LogOut_Response{} }
func (m *User_LogOut_Response) String() string            { return proto.CompactTextString(m) }
func (*User_LogOut_Response) ProtoMessage()               {}
func (*User_LogOut_Response) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

type User_RefreshToken_Request struct {
	RefreshToken string `protobuf:"bytes,1,opt,name=refresh_token,json=refreshToken" json:"refresh_token,omitempty"`
}

func (m *User_RefreshToken_Request) Reset()                    { *m = User_RefreshToken_Request{} }
func (m *User_RefreshToken_Request) String() string            { return proto.CompactTextString(m) }
func (*User_RefreshToken_Request) ProtoMessage()               {}
func (*User_RefreshToken_Request) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

func (m *User_RefreshToken_Request) GetRefreshToken() string {
	if m != nil {
		return m.RefreshToken
	}
	return ""
}

type User_RefreshToken_Response struct {
	TokenData *UserTokenData `protobuf:"bytes,1,opt,name=tokenData" json:"tokenData,omitempty"`
}

func (m *User_RefreshToken_Response) Reset()                    { *m = User_RefreshToken_Response{} }
func (m *User_RefreshToken_Response) String() string            { return proto.CompactTextString(m) }
func (*User_RefreshToken_Response) ProtoMessage()               {}
func (*User_RefreshToken_Response) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

func (m *User_RefreshToken_Response) GetTokenData() *UserTokenData {
	if m != nil {
		return m.TokenData
	}
	return nil
}

func init() {
	proto.RegisterType((*UserTokenData)(nil), "Dolphin.Protocol.UserTokenData")
	proto.RegisterType((*UserData)(nil), "Dolphin.Protocol.UserData")
	proto.RegisterType((*User_Login_Request)(nil), "Dolphin.Protocol.User_Login_Request")
	proto.RegisterType((*User_Login_Response)(nil), "Dolphin.Protocol.User_Login_Response")
	proto.RegisterType((*User_Register_Request)(nil), "Dolphin.Protocol.User_Register_Request")
	proto.RegisterType((*User_Register_Response)(nil), "Dolphin.Protocol.User_Register_Response")
	proto.RegisterType((*User_GetVerifyCode_Request)(nil), "Dolphin.Protocol.User_GetVerifyCode_Request")
	proto.RegisterType((*User_GetVerifyCode_Response)(nil), "Dolphin.Protocol.User_GetVerifyCode_Response")
	proto.RegisterType((*User_GetSelfInfo_Request)(nil), "Dolphin.Protocol.User_GetSelfInfo_Request")
	proto.RegisterType((*User_GetSelfInfo_Response)(nil), "Dolphin.Protocol.User_GetSelfInfo_Response")
	proto.RegisterType((*User_LogOut_Request)(nil), "Dolphin.Protocol.User_LogOut_Request")
	proto.RegisterType((*User_LogOut_Response)(nil), "Dolphin.Protocol.User_LogOut_Response")
	proto.RegisterType((*User_RefreshToken_Request)(nil), "Dolphin.Protocol.User_RefreshToken_Request")
	proto.RegisterType((*User_RefreshToken_Response)(nil), "Dolphin.Protocol.User_RefreshToken_Response")
}

func init() { proto.RegisterFile("user.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 446 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x53, 0x4d, 0x8f, 0xd3, 0x30,
	0x10, 0x55, 0x58, 0xa8, 0xb6, 0xb3, 0x14, 0x56, 0x66, 0x37, 0x0a, 0x41, 0x40, 0x65, 0x2e, 0x7b,
	0xca, 0x22, 0x90, 0xb8, 0xb1, 0xe2, 0x23, 0xd2, 0xaa, 0xd2, 0xaa, 0xa0, 0x14, 0x38, 0x80, 0x44,
	0x94, 0x92, 0x49, 0x1b, 0x91, 0xda, 0xc1, 0x76, 0xf8, 0xf8, 0x05, 0xfc, 0x6d, 0x64, 0xc7, 0x71,
	0x93, 0xaa, 0x07, 0xa4, 0xbd, 0xcd, 0x9b, 0xf1, 0xbc, 0x99, 0x79, 0x33, 0x06, 0x68, 0x24, 0x8a,
	0xa8, 0x16, 0x5c, 0x71, 0x72, 0x1c, 0xf3, 0xaa, 0x5e, 0x97, 0x2c, 0x7a, 0xaf, 0xe1, 0x37, 0x5e,
	0x51, 0x84, 0xc9, 0x47, 0x89, 0xe2, 0x03, 0xff, 0x8e, 0x2c, 0xce, 0x54, 0x46, 0x4e, 0xe0, 0x96,
	0xd2, 0x20, 0xf0, 0xa6, 0xde, 0xd9, 0x38, 0x69, 0x01, 0xa1, 0x70, 0x5b, 0x60, 0x21, 0x50, 0xae,
	0xcd, 0xcb, 0xe0, 0x86, 0x09, 0x0e, 0x7c, 0x24, 0x84, 0x43, 0xfc, 0x5d, 0x97, 0x02, 0x5f, 0xab,
	0xe0, 0x60, 0xea, 0x9d, 0x1d, 0x24, 0x0e, 0xd3, 0x0b, 0x38, 0xd4, 0x65, 0x4c, 0x85, 0xb0, 0xb5,
	0xe7, 0xd9, 0x06, 0x6d, 0x11, 0x87, 0x89, 0x0f, 0x23, 0x6d, 0xcf, 0x62, 0x5b, 0xc1, 0x22, 0x7a,
	0x05, 0x44, 0x5b, 0xe9, 0x15, 0x5f, 0x95, 0x2c, 0x4d, 0xf0, 0x47, 0x83, 0x52, 0x69, 0x26, 0x3d,
	0x1c, 0xeb, 0x31, 0x75, 0x58, 0xc7, 0xea, 0x4c, 0xca, 0x5f, 0x5c, 0xe4, 0x96, 0xcb, 0x61, 0xfa,
	0x15, 0xee, 0x0d, 0xd8, 0x64, 0xcd, 0x99, 0x44, 0x72, 0x09, 0x77, 0x75, 0x7a, 0x6a, 0x46, 0x4e,
	0xf3, 0x4c, 0x65, 0x86, 0xf5, 0xe8, 0xd9, 0xe3, 0x68, 0x57, 0xb7, 0x68, 0x20, 0x5a, 0x32, 0x69,
	0xfa, 0x90, 0xfe, 0xf5, 0xe0, 0xd4, 0x14, 0x48, 0x70, 0x55, 0x4a, 0x65, 0x8c, 0x6b, 0x75, 0x4c,
	0x1e, 0x01, 0xfc, 0x44, 0x51, 0x16, 0x7f, 0xde, 0xf2, 0x1c, 0x8d, 0xba, 0xe3, 0xa4, 0xe7, 0xd1,
	0xba, 0xf1, 0x1a, 0xd9, 0x2c, 0x0f, 0x6e, 0xb6, 0xba, 0xb5, 0x88, 0x3e, 0x05, 0x7f, 0xb7, 0x11,
	0x3b, 0xac, 0x0f, 0x23, 0x81, 0xb2, 0xa9, 0x94, 0xed, 0xc3, 0x22, 0x7a, 0x01, 0xa1, 0xc9, 0xb8,
	0x44, 0xf5, 0xc9, 0xf1, 0xbb, 0xfe, 0xa7, 0x70, 0x54, 0xaf, 0x39, 0xc3, 0x79, 0xb3, 0x59, 0xa2,
	0xb0, 0xa9, 0x7d, 0x17, 0x7d, 0x08, 0x0f, 0xf6, 0xe6, 0xb7, 0x65, 0x69, 0x08, 0x41, 0x17, 0x5e,
	0x60, 0x55, 0xcc, 0x58, 0xc1, 0x3b, 0x72, 0xba, 0x80, 0xfb, 0x7b, 0x62, 0xb6, 0xdf, 0x17, 0xad,
	0x72, 0xf1, 0x76, 0x2b, 0xe1, 0xfe, 0xad, 0x98, 0x85, 0xb8, 0xb7, 0xf4, 0x74, 0xbb, 0xeb, 0x77,
	0x8d, 0x72, 0xb5, 0x7c, 0x38, 0x19, 0xba, 0x6d, 0x7f, 0xaf, 0x6c, 0x0f, 0x49, 0xef, 0xb2, 0xdd,
	0xf4, 0x4f, 0x60, 0x62, 0x2f, 0x3e, 0xed, 0xff, 0x91, 0xc1, 0x37, 0xa0, 0x5f, 0xac, 0x80, 0x3b,
	0x0c, 0x76, 0x8c, 0x97, 0x30, 0x56, 0xdd, 0x9d, 0xfc, 0xef, 0x75, 0x6d, 0x33, 0xde, 0x1c, 0x7f,
	0xbe, 0x13, 0x45, 0xe7, 0x1b, 0x9e, 0x37, 0x15, 0x9e, 0xeb, 0x19, 0x97, 0x23, 0xf3, 0xb3, 0x9f,
	0xff, 0x0b, 0x00, 0x00, 0xff, 0xff, 0xb3, 0x23, 0x54, 0x11, 0xe7, 0x03, 0x00, 0x00,
}
