package model

type HTTPError struct {
	who     string
	query   string
	message string
	code    int
}

func NewHTTPError(who, query, message string, code int) *HTTPError {
	return &HTTPError{who, query, message, code}
}

func (e *HTTPError) Error() string {
	return e.message
}

func (e *HTTPError) Code() int {
	return e.code
}

func (e *HTTPError) Who() string {
	return e.who
}

func (e *HTTPError) Query() string {
	return e.query
}
