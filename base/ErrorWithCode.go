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
	ErrNil              = NewError(200, "")
	ErrBadRequest       = NewError(400, "Bad Request")
	ErrUnauthorized     = NewError(401, "Unauthorized ")
	ErrNotFound         = NewError(404, "Not Found")
	ErrMethodNotAllowed = NewError(405, "Method Not Allowed")
	ErrNotAcceptable    = NewError(406, "Not Acceptable")
	ErrRequestTimeout   = NewError(408, "Request Timeout")
	ErrParamParseFail   = NewError(451, "Param Parse Fail")
	ErrParamNotAllow    = NewError(452, "Param Not Allow")
	ErrInvalidToken     = NewError(498, "Invalid Token")
	ErrRequestToken     = NewError(499, "Token Required")
	ErrLoginTimeOut     = NewError(440, "The client's session has expired and must log in again")
	ErrLoginFail        = NewError(441, "UserName or Password Error")
	ErrServerIsDown     = NewError(521, "Server Is Down")
	ErrUnknown          = NewError(520, "Error Unknown")
	ErrFrozen           = NewError(530, "Frozen")
	ErrInternal         = NewError(500, "Internal Server Error")
)
