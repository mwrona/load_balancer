package model

type HTTPError struct {
	message string
	code    int
}

func NewHTTPError(message string, code int) *HTTPError {
	return &HTTPError{message, code}
}

func (e *HTTPError) Error() string {
	return e.message
}

func (e *HTTPError) Code() int {
	return e.code
}
