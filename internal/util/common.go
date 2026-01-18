package util

import "mini-crypto-exchange/internal/apperrors"

type Error struct {
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}

type RouterConfig struct {
	MatchingEngine interface{}
}

func ServerToError(err error) *Error {
	errorProto := &Error{
		Code:    err.(*apperrors.ServerError).Code,
		Message: err.Error(),
	}
	return errorProto
}
