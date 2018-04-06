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
	ErrNil               = NewError(0x00, "")
	ErrInternal          = NewError(0xff, "Internal Server Error")
	ErrParamParseFail    = NewError(0xfe, "Param Parse Fail")
	ErrParamNotAllow     = NewError(0xfd, "Param Not Allow")
	ErrBadRequest        = NewError(0xfa, "Bad Request")
	ErrUnauthorized      = NewError(0xfb, "Unauthorized ")
	ErrNotFound          = NewError(0xfc, "Not Found")
	ErrNeedLogin         = NewError(0xf9, "Need Login first")
	ErrLoginTimeOut      = NewError(0xf8, "The client's session has expired and must log in again")
	ErrServerIsDown      = NewError(0xf1, "Server Is Down")
	ErrLoginFail         = NewError(1, "UserName or Password Error")
	ErrAccountBeenTaken  = NewError(2, "Account has been taken")
	ErrNameOrPwdShort    = NewError(1, "Name or Password is too short")
	ErrVerifyCodeErr     = NewError(3, "Verify Code Error ")
	ErrVerifySendTooBusy = NewError(1, "Verify Send too busy ")
	ErrUnknown           = NewError(0xf7, "Error Unknown")
	ErrFrozen            = NewError(0xf6, "Frozen")
	ErrMethodNotAllowed  = NewError(0xf5, "Method Not Allowed")
	ErrNotAcceptable     = NewError(0xf4, "Not Acceptable")
	ErrRequestTimeout    = NewError(0xf3, "Request Timeout")
	ErrInvalidToken      = NewError(0xf2, "Invalid Token")
	ErrRequestToken      = NewError(0xef, "Token Required")
	ErrSQLERROR          = NewError(0xee, "SQL ERROR")
	ErrSMSSendFail       = NewError(0xef, "SMS Send Fail")
)
