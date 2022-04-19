package http

import (
	"fmt"
)

type Error struct {
	Code    uint64
	Message string
}

func (error *Error) Serialize() string {
	return fmt.Sprintf(`{"error": %d, "message": "%s"}`, error.Code, error.Message)
}

func ToReadable(error *Error) string {
	if error == nil {
		return "SUCCESS"
	}
	return fmt.Sprintf(`%d %s`, error.Code, error.Message)
}

const (
	YES uint64 = 1
	NO  uint64 = 0
)

const (
	INVALID_REQUEST        uint64 = 101
	INTERNAL_SERVER_ERROR  uint64 = 102
	AUTHORIZATION_FAILED   uint64 = 103
	INVALID_CURRENCY_CODE  uint64 = 104
	CURRENCY_CODE_MISMATCH uint64 = 105
	NOT_FOUND              uint64 = 106
	WRONG_AMOUNT           uint64 = 107
	INSUFFICIENT_FUNDS     uint64 = 108
	TOO_MANY_REQUESTS      uint64 = 109
	ALREADY_EXISTS         uint64 = 110
)

func ErrorCodeToString(code uint64) string {
	if code == INVALID_REQUEST {
		return "INVALID REQUEST"
	}
	if code == INTERNAL_SERVER_ERROR {
		return "INTERNAL SERVER ERROR"
	}
	if code == AUTHORIZATION_FAILED {
		return "AUTHORIZATION FAILED"
	}
	if code == INVALID_CURRENCY_CODE {
		return "INVALID CURRENCY CODE"
	}
	if code == CURRENCY_CODE_MISMATCH {
		return "CURRENCY CODE MISMATCH"
	}
	if code == WRONG_AMOUNT {
		return "WRONG AMOUNT"
	}
	if code == NOT_FOUND {
		return "NOT FOUND"
	}
	if code == INSUFFICIENT_FUNDS {
		return "INSUFFICIENT FUNDS"
	}
	if code == TOO_MANY_REQUESTS {
		return "TOO MANY REQUESTS"
	}
	return fmt.Sprintf("%v", code)
}

func ToError(code uint64) *Error {
	return &Error{
		Code:    code,
		Message: ErrorCodeToString(code),
	}
}
