package helper

type APIError struct {
	Status  int
	Code    string
	Message string
	Cause   error
	Fields  map[string]string
}

func (e *APIError) Error() string {
	return e.Message
}

func NewAPIError(status int, code, message string) *APIError {
	return &APIError{
		Status:  status,
		Code:    code,
		Message: message,
	}
}

func WithCause(e *APIError, cause error) *APIError {
	e.Cause = cause
	return e
}

func Validation(fields map[string]string) *APIError {
	return &APIError{
		Status:  422,
		Code:    "VALIDATION_ERROR",
		Message: "validasi gagal",
		Fields:  fields,
	}
}

func BadRequest(m string) *APIError {
	return NewAPIError(400, "BAD_REQUEST", m)
}

func Unauthorized(m string) *APIError {
	return NewAPIError(401, "UNAUTHORIZED", m)
}

func Forbidden(m string) *APIError {
	return NewAPIError(403, "FORBIDDEN", m)
}

func NotFound(m string) *APIError {
	return NewAPIError(404, "NOT_FOUND", m)
}

func Conflict(m string) *APIError {
	return NewAPIError(409, "CONFLICT", m)
}

func Unsupported(m string) *APIError {
	return NewAPIError(415, "UNSUPPORTED_MEDIA_TYPE", m)
}

func NotAcceptable(m string) *APIError {
	return NewAPIError(406, "NOT_ACCEPTABLE", m)
}

func Internal(cause error) *APIError {
	return WithCause(
		NewAPIError(
			500,
			"INTERNAL_ERROR",
			"terjadi kesalahan pada server",
		),
		cause,
	)
}

func BadID() *APIError {
	return BadRequest("id tidak valid")
}
