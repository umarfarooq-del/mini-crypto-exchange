package apperrors

import (
	"net/http"
)

type ServerError struct {
	Message          string
	Code             string
	HTTPResponseCode int
	GRPCResponseCode uint32
}

func (err *ServerError) Error() string {
	return err.Code + ": " + err.Message
}

var (
	ErrPairNotFound = &ServerError{
		Code:             "PAIR_NOT_FOUND",
		Message:          "Trading pair not found",
		HTTPResponseCode: http.StatusNotFound,
	}

	ErrInvalidPrice = &ServerError{
		Code:             "INVALID_PRICE",
		Message:          "Price must be greater than 0",
		HTTPResponseCode: http.StatusBadRequest,
	}

	ErrInvalidQuantity = &ServerError{
		Code:             "INVALID_QUANTITY",
		Message:          "Quantity must be greater than 0",
		HTTPResponseCode: http.StatusBadRequest,
	}

	ErrInvalidSide = &ServerError{
		Code:             "INVALID_SIDE",
		Message:          "Side must be 'buy' or 'sell'",
		HTTPResponseCode: http.StatusBadRequest,
	}

	ErrInvalidUserID = &ServerError{
		Code:             "INVALID_USER_ID",
		Message:          "User ID must be greater than 0",
		HTTPResponseCode: http.StatusBadRequest,
	}
	
)
