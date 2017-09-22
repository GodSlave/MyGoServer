package base

import "fmt"

type ErrorCode struct {
	ErrorCode int32
	Desc      string
}

func (e ErrorCode) Error() string {
	return fmt.Sprintf("%v: %v", e.ErrorCode, e.Desc)
}

func NewError(errorCode int32, text string) *ErrorCode {
	return &ErrorCode{errorCode, text}
}

var (
	ErrNil              = NewError(0x00, "")
	ErrInternal         = NewError(0xff, "Internal Server Error")
	ErrParamParseFail   = NewError(0xfe, "Param Parse Fail")
	ErrParamNotAllow    = NewError(0xfd, "Param Not Allow")
	ErrBadRequest       = NewError(0xfa, "Bad Request")
	ErrUnauthorized     = NewError(0xfb, "Unauthorized ")
	ErrNotFound         = NewError(0xfc, "Not Found")
	ErrNeedLogin        = NewError(0xf9, "Need Login first")
	ErrLoginTimeOut     = NewError(0xf8, "The client's session has expired and must log in again")
	ErrServerIsDown     = NewError(0xf1, "Server Is Down")
	ErrLoginFail        = NewError(0x020101, "UserName or Password Error")
	ErrAccountBeenTaken = NewError(0x020201, "Account has been taken")
	ErrNameOrPwdShort   = NewError(0x020202, "Name or Password is too short")
	ErrVerifyCodeErr    = NewError(0x020203, "Verify Code Error ")
	ErrUnknown          = NewError(520, "Error Unknown")
	ErrFrozen           = NewError(530, "Frozen")
	ErrMethodNotAllowed = NewError(405, "Method Not Allowed")
	ErrNotAcceptable    = NewError(406, "Not Acceptable")
	ErrRequestTimeout   = NewError(408, "Request Timeout")
	ErrInvalidToken     = NewError(498, "Invalid Token")
	ErrRequestToken     = NewError(499, "Token Required")
)
