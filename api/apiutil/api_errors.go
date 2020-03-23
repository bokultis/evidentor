package apiutil

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Error is encapsulated error fulfilling GO error interface
type Error struct {
	Code     string
	Message  string
	HTTPCode int
	Items    []*ErrorItem
}

type ErrorItem struct {
	Name    string `json:"name"`
	Message string `json:"message,omitempty"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("CODE: %s,MSG: %s,HTTP STS: %d, ITEMS: %s", e.Code, e.Message, e.HTTPCode, e.Items)
}

func (e *ErrorItem) String() string {
	return fmt.Sprintf("ITEM=[%s:%s]", e.Name, e.Message)
}

// NewError creates new API error
func NewError(httpCode int, code string, message string) *Error {
	return &Error{
		Code:     code,
		Message:  message,
		HTTPCode: httpCode,
	}
}

// NewErrorItem is ErrorItem constructor
func NewErrorItem(name, message string) *ErrorItem {
	return &ErrorItem{
		Name:    name,
		Message: message,
	}
}

// NewIntError is internal server error
func NewIntError(e error) *Error {
	return &Error{
		Code:     "INT_APP_ERROR",
		Message:  e.Error(),
		HTTPCode: http.StatusInternalServerError,
	}
}

// NewValidationError is validation data error
func NewValidationError(items []*ErrorItem) *Error {
	return &Error{
		Code:     "INVALID-DATA",
		Message:  "Input data validation failed",
		HTTPCode: http.StatusBadRequest,
		Items:    items,
	}
}

var (
	// ErrNotAuthenticated is returned when call is not authenticated
	ErrNotAuthenticated = NewError(http.StatusUnauthorized,
		"NOT_AUTHENTICATED", "Not authenticated")

	// ErrFailedEncoding when JSON encoding/decoding fails
	ErrFailedEncoding = NewError(http.StatusInternalServerError,
		"FAILED_ENCODING", "Failed encoding")
	// ErrBadParameter is returned when parameter has bad value
	ErrBadParameter = NewError(http.StatusNotFound,
		"BAD_PARAM", "Bad parameter")

	ErrBadPageNumber = NewError(http.StatusBadRequest,
		"BAD_PARAM_PAGE_NUMBER", "Bad parameter 'page_number': value must be greather than 0")
	//ErrBadPageSize is error returned when page_size is less than 1
	ErrBadPageSize = NewError(http.StatusBadRequest,
		"BAD_PARAM_PAGE_SIZE", "Bad parameter 'page_size': value : must be greather than 0")
)

//ErrUnsupportedFilterProperty is returned when filter property with given name is not supported
var ErrUnsupportedFilterProperty = NewError(http.StatusBadRequest,
	"REQUEST_ERROR", "unsupported filter property")

func NewErrorResponse(w http.ResponseWriter, response *Error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.HTTPCode)
	json.NewEncoder(w).Encode(response)
	return
}
