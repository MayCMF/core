package errors

// ResponseError - Define response error
type ResponseError struct {
	Code       int    // error code
	Message    string // wrong information
	StatusCode int    // Response status code
	ERR        error  // Response error
}

func (r *ResponseError) Error() string {
	if r.ERR != nil {
		return r.ERR.Error()
	}
	return r.Message
}

// UnWrapResponse - Unpacking response error
func UnWrapResponse(err error) *ResponseError {
	if v, ok := err.(*ResponseError); ok {
		return v
	}
	return nil
}

// WrapResponse - Wrapper response error
func WrapResponse(err error, code int, msg string, status ...int) error {
	res := &ResponseError{
		Code:    code,
		Message: msg,
		ERR:     err,
	}
	if len(status) > 0 {
		res.StatusCode = status[0]
	}
	return res
}

// Wrap400Response - Wrong response error with package error code 400
func Wrap400Response(err error, msg ...string) error {
	m := "Request error"
	if len(msg) > 0 {
		m = msg[0]
	}
	return WrapResponse(err, 400, m, 400)
}

// Wrap500Response - Wrong response error with package error code 500
func Wrap500Response(err error, msg ...string) error {
	m := "Server error"
	if len(msg) > 0 {
		m = msg[0]
	}
	return WrapResponse(err, 500, m, 500)
}

// NewResponse - Create response error
func NewResponse(code int, msg string, status ...int) error {
	res := &ResponseError{
		Code:    code,
		Message: msg,
	}
	if len(status) > 0 {
		res.StatusCode = status[0]
	}
	return res
}

// New400Response - Create a response error with error code 400
func New400Response(msg string) error {
	return NewResponse(400, msg, 400)
}

// New500Response - Create a response error with error code 500
func New500Response(msg string) error {
	return NewResponse(500, msg, 500)
}
