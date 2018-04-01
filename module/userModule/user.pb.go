// Code generated by protoc-gen-go. DO NOT EDIT.
// source: user.proto

/*
Package userModule is a generated protocol buffer package.

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
	User_ForgetPassWord_Request
	User_ForgetPassWord_Response
	User_ChangePassWordByVerifyCode_Request
	User_ChangePassWordByVerifyCode_Response
	User_ChangePassWordByPassword_Request
	User_ChangePassWordByPassword_Response
	User_CheckVerifyCodeAvailable_Request
	User_CheckVerifyCodeAvailable_Response
*/
package userModule

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
	Result    string         `protobuf:"bytes,1,opt,name=result" json:"result,omitempty"`
	TokenInfo *UserTokenData `protobuf:"bytes,2,opt,name=tokenInfo" json:"tokenInfo,omitempty"`
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

func (m *User_Register_Response) GetTokenInfo() *UserTokenData {
	if m != nil {
		return m.TokenInfo
	}
	return nil
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

type User_ForgetPassWord_Request struct {
	PhoneNumber string `protobuf:"bytes,1,opt,name=phoneNumber" json:"phoneNumber,omitempty"`
}

func (m *User_ForgetPassWord_Request) Reset()                    { *m = User_ForgetPassWord_Request{} }
func (m *User_ForgetPassWord_Request) String() string            { return proto.CompactTextString(m) }
func (*User_ForgetPassWord_Request) ProtoMessage()               {}
func (*User_ForgetPassWord_Request) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{14} }

func (m *User_ForgetPassWord_Request) GetPhoneNumber() string {
	if m != nil {
		return m.PhoneNumber
	}
	return ""
}

type User_ForgetPassWord_Response struct {
}

func (m *User_ForgetPassWord_Response) Reset()                    { *m = User_ForgetPassWord_Response{} }
func (m *User_ForgetPassWord_Response) String() string            { return proto.CompactTextString(m) }
func (*User_ForgetPassWord_Response) ProtoMessage()               {}
func (*User_ForgetPassWord_Response) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{15} }

type User_ChangePassWordByVerifyCode_Request struct {
	PhoneNumber string `protobuf:"bytes,1,opt,name=phoneNumber" json:"phoneNumber,omitempty"`
	VerifyCode  string `protobuf:"bytes,2,opt,name=verifyCode" json:"verifyCode,omitempty"`
	Password    string `protobuf:"bytes,3,opt,name=password" json:"password,omitempty"`
}

func (m *User_ChangePassWordByVerifyCode_Request) Reset() {
	*m = User_ChangePassWordByVerifyCode_Request{}
}
func (m *User_ChangePassWordByVerifyCode_Request) String() string { return proto.CompactTextString(m) }
func (*User_ChangePassWordByVerifyCode_Request) ProtoMessage()    {}
func (*User_ChangePassWordByVerifyCode_Request) Descriptor() ([]byte, []int) {
	return fileDescriptor0, []int{16}
}

func (m *User_ChangePassWordByVerifyCode_Request) GetPhoneNumber() string {
	if m != nil {
		return m.PhoneNumber
	}
	return ""
}

func (m *User_ChangePassWordByVerifyCode_Request) GetVerifyCode() string {
	if m != nil {
		return m.VerifyCode
	}
	return ""
}

func (m *User_ChangePassWordByVerifyCode_Request) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type User_ChangePassWordByVerifyCode_Response struct {
	IsSuccess bool `protobuf:"varint,1,opt,name=isSuccess" json:"isSuccess,omitempty"`
}

func (m *User_ChangePassWordByVerifyCode_Response) Reset() {
	*m = User_ChangePassWordByVerifyCode_Response{}
}
func (m *User_ChangePassWordByVerifyCode_Response) String() string { return proto.CompactTextString(m) }
func (*User_ChangePassWordByVerifyCode_Response) ProtoMessage()    {}
func (*User_ChangePassWordByVerifyCode_Response) Descriptor() ([]byte, []int) {
	return fileDescriptor0, []int{17}
}

func (m *User_ChangePassWordByVerifyCode_Response) GetIsSuccess() bool {
	if m != nil {
		return m.IsSuccess
	}
	return false
}

type User_ChangePassWordByPassword_Request struct {
	OldPassword string `protobuf:"bytes,1,opt,name=oldPassword" json:"oldPassword,omitempty"`
	NewPassword string `protobuf:"bytes,2,opt,name=newPassword" json:"newPassword,omitempty"`
}

func (m *User_ChangePassWordByPassword_Request) Reset()         { *m = User_ChangePassWordByPassword_Request{} }
func (m *User_ChangePassWordByPassword_Request) String() string { return proto.CompactTextString(m) }
func (*User_ChangePassWordByPassword_Request) ProtoMessage()    {}
func (*User_ChangePassWordByPassword_Request) Descriptor() ([]byte, []int) {
	return fileDescriptor0, []int{18}
}

func (m *User_ChangePassWordByPassword_Request) GetOldPassword() string {
	if m != nil {
		return m.OldPassword
	}
	return ""
}

func (m *User_ChangePassWordByPassword_Request) GetNewPassword() string {
	if m != nil {
		return m.NewPassword
	}
	return ""
}

type User_ChangePassWordByPassword_Response struct {
	IsSuccess bool `protobuf:"varint,1,opt,name=isSuccess" json:"isSuccess,omitempty"`
}

func (m *User_ChangePassWordByPassword_Response) Reset() {
	*m = User_ChangePassWordByPassword_Response{}
}
func (m *User_ChangePassWordByPassword_Response) String() string { return proto.CompactTextString(m) }
func (*User_ChangePassWordByPassword_Response) ProtoMessage()    {}
func (*User_ChangePassWordByPassword_Response) Descriptor() ([]byte, []int) {
	return fileDescriptor0, []int{19}
}

func (m *User_ChangePassWordByPassword_Response) GetIsSuccess() bool {
	if m != nil {
		return m.IsSuccess
	}
	return false
}

type User_CheckVerifyCodeAvailable_Request struct {
	PhoneNumber string `protobuf:"bytes,1,opt,name=phoneNumber" json:"phoneNumber,omitempty"`
	VerifyCode  string `protobuf:"bytes,2,opt,name=verifyCode" json:"verifyCode,omitempty"`
}

func (m *User_CheckVerifyCodeAvailable_Request) Reset()         { *m = User_CheckVerifyCodeAvailable_Request{} }
func (m *User_CheckVerifyCodeAvailable_Request) String() string { return proto.CompactTextString(m) }
func (*User_CheckVerifyCodeAvailable_Request) ProtoMessage()    {}
func (*User_CheckVerifyCodeAvailable_Request) Descriptor() ([]byte, []int) {
	return fileDescriptor0, []int{20}
}

func (m *User_CheckVerifyCodeAvailable_Request) GetPhoneNumber() string {
	if m != nil {
		return m.PhoneNumber
	}
	return ""
}

func (m *User_CheckVerifyCodeAvailable_Request) GetVerifyCode() string {
	if m != nil {
		return m.VerifyCode
	}
	return ""
}

type User_CheckVerifyCodeAvailable_Response struct {
	IsAvalible     bool `protobuf:"varint,1,opt,name=isAvalible" json:"isAvalible,omitempty"`
	IsAccountTaken bool `protobuf:"varint,2,opt,name=isAccountTaken" json:"isAccountTaken,omitempty"`
}

func (m *User_CheckVerifyCodeAvailable_Response) Reset() {
	*m = User_CheckVerifyCodeAvailable_Response{}
}
func (m *User_CheckVerifyCodeAvailable_Response) String() string { return proto.CompactTextString(m) }
func (*User_CheckVerifyCodeAvailable_Response) ProtoMessage()    {}
func (*User_CheckVerifyCodeAvailable_Response) Descriptor() ([]byte, []int) {
	return fileDescriptor0, []int{21}
}

func (m *User_CheckVerifyCodeAvailable_Response) GetIsAvalible() bool {
	if m != nil {
		return m.IsAvalible
	}
	return false
}

func (m *User_CheckVerifyCodeAvailable_Response) GetIsAccountTaken() bool {
	if m != nil {
		return m.IsAccountTaken
	}
	return false
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
	proto.RegisterType((*User_ForgetPassWord_Request)(nil), "Dolphin.Protocol.User_ForgetPassWord_Request")
	proto.RegisterType((*User_ForgetPassWord_Response)(nil), "Dolphin.Protocol.User_ForgetPassWord_Response")
	proto.RegisterType((*User_ChangePassWordByVerifyCode_Request)(nil), "Dolphin.Protocol.User_ChangePassWordByVerifyCode_Request")
	proto.RegisterType((*User_ChangePassWordByVerifyCode_Response)(nil), "Dolphin.Protocol.User_ChangePassWordByVerifyCode_Response")
	proto.RegisterType((*User_ChangePassWordByPassword_Request)(nil), "Dolphin.Protocol.User_ChangePassWordByPassword_Request")
	proto.RegisterType((*User_ChangePassWordByPassword_Response)(nil), "Dolphin.Protocol.User_ChangePassWordByPassword_Response")
	proto.RegisterType((*User_CheckVerifyCodeAvailable_Request)(nil), "Dolphin.Protocol.User_CheckVerifyCodeAvailable_Request")
	proto.RegisterType((*User_CheckVerifyCodeAvailable_Response)(nil), "Dolphin.Protocol.User_CheckVerifyCodeAvailable_Response")
}

func init() { proto.RegisterFile("user.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 626 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x55, 0x6f, 0x4f, 0x13, 0x4f,
	0x10, 0x4e, 0xe1, 0xf7, 0x23, 0xed, 0x20, 0x6a, 0x4e, 0x68, 0x6a, 0x45, 0x24, 0x6b, 0x44, 0x5e,
	0x1d, 0x89, 0x26, 0xbe, 0x13, 0x2d, 0x34, 0x20, 0x09, 0x62, 0x73, 0x45, 0x4d, 0x34, 0xb1, 0xd9,
	0xde, 0x4d, 0xdb, 0x4b, 0xaf, 0xb7, 0xe7, 0xee, 0x5e, 0x91, 0x4f, 0xc0, 0xd7, 0x36, 0xbb, 0xb7,
	0xf7, 0xd7, 0x1a, 0x40, 0x7d, 0xb7, 0xcf, 0xec, 0xee, 0x3c, 0xcf, 0x33, 0x3b, 0x93, 0x05, 0x88,
	0x05, 0x72, 0x3b, 0xe2, 0x4c, 0x32, 0xeb, 0x7e, 0x97, 0x05, 0xd1, 0xc4, 0x0f, 0xed, 0x9e, 0x82,
	0x2e, 0x0b, 0x08, 0xc2, 0xda, 0x47, 0x81, 0xfc, 0x9c, 0x4d, 0x31, 0xec, 0x52, 0x49, 0xad, 0x75,
	0xf8, 0x5f, 0x2a, 0xd0, 0xaa, 0x6d, 0xd7, 0x76, 0x1b, 0x4e, 0x02, 0x2c, 0x02, 0x77, 0x38, 0x8e,
	0x38, 0x8a, 0x89, 0x3e, 0xd9, 0x5a, 0xd2, 0x9b, 0xa5, 0x98, 0xd5, 0x86, 0x3a, 0xfe, 0x88, 0x7c,
	0x8e, 0x1d, 0xd9, 0x5a, 0xde, 0xae, 0xed, 0x2e, 0x3b, 0x19, 0x26, 0xfb, 0x50, 0x57, 0x34, 0x9a,
	0xa1, 0x9d, 0xac, 0xcf, 0xe8, 0x0c, 0x0d, 0x49, 0x86, 0xad, 0x26, 0xac, 0xa8, 0xf5, 0x49, 0xd7,
	0x30, 0x18, 0x44, 0x4e, 0xc1, 0x52, 0xab, 0xc1, 0x29, 0x1b, 0xfb, 0xe1, 0xc0, 0xc1, 0xef, 0x31,
	0x0a, 0xa9, 0x32, 0x29, 0x73, 0x61, 0x21, 0x53, 0x8a, 0xd5, 0x5e, 0x44, 0x85, 0xb8, 0x60, 0xdc,
	0x33, 0xb9, 0x32, 0x4c, 0xbe, 0xc1, 0x83, 0x52, 0x36, 0x11, 0xb1, 0x50, 0xa0, 0x75, 0x0c, 0xf7,
	0xd4, 0xf5, 0x81, 0xb6, 0x3c, 0xf0, 0xa8, 0xa4, 0x3a, 0xeb, 0xea, 0x8b, 0x27, 0x76, 0xb5, 0x6e,
	0x76, 0xa9, 0x68, 0xce, 0x5a, 0x5c, 0x84, 0xe4, 0xaa, 0x06, 0x1b, 0x9a, 0xc0, 0xc1, 0xb1, 0x2f,
	0xa4, 0x5e, 0xfc, 0x95, 0x62, 0x6b, 0x0b, 0x60, 0x8e, 0xdc, 0x1f, 0x5d, 0x1e, 0x32, 0x0f, 0x75,
	0x75, 0x1b, 0x4e, 0x21, 0xa2, 0xea, 0xc6, 0x22, 0x0c, 0x4f, 0xbc, 0xd6, 0x7f, 0x49, 0xdd, 0x12,
	0x44, 0x18, 0x34, 0xab, 0x42, 0x8c, 0xd9, 0x26, 0xac, 0x70, 0x14, 0x71, 0x20, 0x8d, 0x0e, 0x83,
	0xac, 0xd7, 0xd0, 0xd0, 0xfe, 0x4f, 0xc2, 0x11, 0xd3, 0x32, 0x6e, 0x60, 0x3f, 0xbf, 0x41, 0xf6,
	0xa1, 0xad, 0x09, 0x8f, 0x51, 0x7e, 0xca, 0xe4, 0x65, 0xf6, 0xb7, 0x61, 0x35, 0x9a, 0xb0, 0x10,
	0xcf, 0xe2, 0xd9, 0x10, 0xb9, 0x61, 0x2e, 0x86, 0xc8, 0x63, 0x78, 0xb4, 0xf0, 0x7e, 0xa2, 0x9a,
	0xb4, 0xa1, 0x95, 0x6e, 0xf7, 0x31, 0x18, 0x29, 0xca, 0x34, 0x39, 0xe9, 0xc3, 0xc3, 0x05, 0x7b,
	0xc6, 0xee, 0xab, 0xa4, 0xf0, 0xdd, 0xfc, 0x51, 0xdb, 0x8b, 0x5d, 0x69, 0x43, 0xd9, 0x59, 0xb2,
	0x91, 0xb7, 0xca, 0x87, 0x58, 0x66, 0x5c, 0x4d, 0x58, 0x2f, 0x87, 0x8d, 0xbe, 0xb7, 0x46, 0x83,
	0x53, 0x18, 0x8c, 0xcc, 0xfd, 0x53, 0x58, 0x33, 0x03, 0x33, 0x28, 0x8e, 0x58, 0x69, 0x8a, 0xc8,
	0x57, 0x53, 0xc0, 0x4a, 0x06, 0x63, 0x23, 0x7d, 0x9d, 0xee, 0x2d, 0x9a, 0x33, 0xbf, 0x41, 0xde,
	0x98, 0xea, 0x1e, 0x31, 0x3e, 0x46, 0xd9, 0xa3, 0x42, 0x7c, 0x66, 0xdc, 0xbb, 0xc5, 0xf3, 0x6c,
	0xc1, 0xe6, 0xe2, 0x04, 0xc6, 0xff, 0x55, 0x0d, 0x9e, 0xeb, 0x03, 0x87, 0x13, 0x1a, 0x8e, 0x31,
	0x3d, 0x70, 0x70, 0xf9, 0x27, 0xcd, 0x50, 0xe9, 0xfa, 0xa5, 0x5f, 0xba, 0xbe, 0x38, 0x31, 0xcb,
	0x95, 0x19, 0x7f, 0x07, 0xbb, 0xd7, 0x0b, 0x31, 0x55, 0xdd, 0x84, 0x86, 0x2f, 0xfa, 0xb1, 0xeb,
	0xa2, 0x10, 0x5a, 0x47, 0xdd, 0xc9, 0x03, 0x64, 0x0a, 0xcf, 0x16, 0x66, 0xea, 0x19, 0xaa, 0xa2,
	0x21, 0x16, 0x78, 0x69, 0x38, 0x35, 0x54, 0x08, 0xa9, 0x13, 0x21, 0x5e, 0xf4, 0xca, 0x53, 0x5e,
	0x0c, 0x91, 0x23, 0xd8, 0xb9, 0x8e, 0xec, 0x46, 0xa2, 0xfd, 0x4c, 0x34, 0xba, 0xd3, 0xdc, 0x73,
	0x67, 0x4e, 0xfd, 0x80, 0x0e, 0x83, 0x7f, 0xf8, 0x0a, 0x24, 0xca, 0x24, 0xff, 0x96, 0xca, 0x48,
	0xde, 0x02, 0xf0, 0x45, 0x67, 0x4e, 0x03, 0x7f, 0x18, 0xa0, 0xd1, 0x5c, 0x88, 0x58, 0x3b, 0x70,
	0xd7, 0x17, 0x1d, 0xd7, 0x65, 0x71, 0x28, 0xcf, 0x69, 0xfa, 0xcf, 0xd4, 0x9d, 0x4a, 0xf4, 0xa0,
	0xf9, 0x65, 0xdd, 0xb6, 0xf7, 0x66, 0xcc, 0x8b, 0x03, 0xdc, 0x53, 0xa3, 0xfa, 0x5e, 0x2f, 0x87,
	0x2b, 0xfa, 0x97, 0x7b, 0xf9, 0x33, 0x00, 0x00, 0xff, 0xff, 0x16, 0x07, 0x8c, 0x56, 0xf3, 0x06,
	0x00, 0x00,
}
