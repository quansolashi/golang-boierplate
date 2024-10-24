package util

import (
	"context"
	"errors"
	"net/http"
)

// ErrorResponse - エラーレスポンス.
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

func NewErrorResponse(err error) (*ErrorResponse, int) {
	if status, ok := internalError(err); ok {
		return newErrorResponse(status, err), status
	}

	if err == nil {
		err = ErrUnknown
	}

	res := &ErrorResponse{
		Status:  http.StatusInternalServerError,
		Message: "Unknown Error",
		Detail:  err.Error(),
	}
	return res, http.StatusInternalServerError
}

func newErrorResponse(status int, err error) *ErrorResponse {
	return &ErrorResponse{
		Status:  status,
		Message: http.StatusText(status),
		Detail:  err.Error(),
	}
}

func (r *ErrorResponse) GetDetail() string {
	if r == nil {
		return ""
	}
	return r.Detail
}

func internalError(err error) (int, bool) {
	if err == nil {
		return 0, false
	}

	var s int
	switch {
	// 4xx
	case // 400
		errors.Is(err, ErrInvalidArgument):
		s = http.StatusBadRequest
	case // 401
		errors.Is(err, ErrUnauthenticated):
		s = http.StatusUnauthorized
	case // 403
		errors.Is(err, ErrForbidden):
		s = http.StatusForbidden
	case // 404
		errors.Is(err, ErrNotFound):
		s = http.StatusNotFound
	case // 406
		errors.Is(err, ErrNotAcceptable):
		s = http.StatusNotAcceptable
	case // 409
		errors.Is(err, ErrAlreadyExists):
		s = http.StatusConflict
	case // 412
		errors.Is(err, ErrFailedPrecondition):
		s = http.StatusPreconditionFailed
	// 5xx
	case // 500
		errors.Is(err, ErrInternal):
		s = http.StatusInternalServerError
	case // 501
		errors.Is(err, ErrNotImplemented):
		s = http.StatusNotImplemented
	case // 502
		errors.Is(err, ErrUnavailable):
		s = http.StatusBadGateway
	case // 504
		errors.Is(err, ErrDeadlineExceeded),
		errors.Is(err, context.Canceled),
		errors.Is(err, context.DeadlineExceeded):
		s = http.StatusGatewayTimeout
	default:
		return 0, false
	}

	return s, true
}
