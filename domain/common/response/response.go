package response

import (
	"encoding/json"
	"net/http"
)

const (
	STATUS_SUCCESS   = "success"
	STATUS_NOT_FOUND = "not found"

	STATUS_REDIRECT    = "redirect"
	STATUS_ERROR       = "internal server error"
	STATUS_BAD_REQUEST = "bad request"

	ERROR_WALLET_DISABLED       = "wallet disabled"
	ERROR_WALLET_NOT_FOUND      = "wallet not found"
	ERROR_INSSUFICIENT_FUND     = "insufficient fund"
	ERROR_REFERENCE_ID_CONFLICT = "reference id already used"
	ERROR_UNAUTHORIZED          = "unauthorized"
)

type Error struct {
	Error string `json:"error"`
}

type Response[T any] struct {
	Status     string              `json:"status"`
	Data       *T                  `json:"data,omitempty"`
	Message    *string             `json:"message,omitempty"`
	StatusCode int                 `json:"-"`
	Writer     http.ResponseWriter `json:"-"`
	Cookies    []*http.Cookie      `json:"-"`
}

func (res *Response[T]) Success(data T) {
	res.Status = STATUS_SUCCESS
	res.StatusCode = http.StatusOK
	res.Data = &data
}

func (res *Response[T]) SuccessWithMessage(message string) {
	res.Status = STATUS_SUCCESS
	res.StatusCode = http.StatusOK
	res.Message = &message
}

func (res *Response[T]) SuccessWithCookie(msg string, data T, cookies []*http.Cookie) {
	res.Status = STATUS_SUCCESS
	res.StatusCode = http.StatusOK
	res.Data = &data
	res.Cookies = cookies
}

func (res *Response[T]) BadRequest(msg string, data *T) {
	res.Status = STATUS_BAD_REQUEST
	res.StatusCode = http.StatusBadRequest
	res.Message = &msg

	if data != nil {
		res.Data = data
	}
}

func (res *Response[T]) NotFound(msg string, data *T) {
	res.Status = STATUS_NOT_FOUND
	res.StatusCode = http.StatusNotFound
	res.Message = &msg

	if data != nil {
		res.Data = data
	}
}

func (res *Response[T]) InternalServerError(msg string) {
	res.Status = STATUS_ERROR
	res.StatusCode = http.StatusInternalServerError
}

func (res *Response[T]) Unauthorized(msg string) {
	res.Status = ERROR_UNAUTHORIZED
	res.StatusCode = http.StatusUnauthorized
}

func (res *Response[T]) Redirect(msg string) {
	res.Status = STATUS_REDIRECT
	res.StatusCode = http.StatusTemporaryRedirect
}

func (res *Response[T]) WriteResponse() {
	if len(res.Cookies) > 0 {
		for _, cookie := range res.Cookies {
			http.SetCookie(res.Writer, cookie)
		}
	}

	res.Writer.Header().Set("Content-Type", "application/json")
	res.Writer.WriteHeader(res.StatusCode)
	json.NewEncoder(res.Writer).Encode(res)
}
