package model

type AppError struct {
	Err          error  `json:"-"`
	StatusCode   int    `json:"-"`
	ErrorCode    string `json:"code"`
	ErrorMessage string `json:"message"`
}

func (e *AppError) Error() string {
	return e.ErrorMessage
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewErrBadResuest(code string, msg string) *AppError {
	return &AppError{StatusCode: 400, ErrorCode: code, ErrorMessage: msg}
}

func NewErrUnauthorized(code string, msg string) *AppError {
	return &AppError{StatusCode: 401, ErrorCode: code, ErrorMessage: msg}
}

func NewErrForbidden(code string, msg string) *AppError {
	return &AppError{StatusCode: 404, ErrorCode: code, ErrorMessage: msg}
}

func NewErrNotFound(code string, msg string) *AppError {
	return &AppError{StatusCode: 404, ErrorCode: code, ErrorMessage: msg}
}

func NewErrUnprocessable(code string, msg string) *AppError {
	return &AppError{StatusCode: 422, ErrorCode: code, ErrorMessage: msg}
}

func NewErrInternalServerError(code string, msg string) *AppError {
	return &AppError{StatusCode: 500, ErrorCode: code, ErrorMessage: msg}
}

func NewAppError(scode int, ecode string, msg string) *AppError {
	return &AppError{StatusCode: scode, ErrorCode: ecode, ErrorMessage: msg}
}
