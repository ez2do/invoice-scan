package errors

var (
	ErrDataNotFound   = &ErrorCode{message: "data not found"}
	ErrDuplicateEntry = &ErrorCode{message: "create duplicate entry"}
)

// ErrorCode
type ErrorCode struct {
	code    string
	message string
}

func (e *ErrorCode) Error() string {
	return e.message
}

func (e *ErrorCode) Code() string {
	return e.code
}

// New create error with message
func New(msg string) error {
	return &ErrorCode{message: msg}
}

// NewCode create error with code, message
func NewCode(code, msg string) error {
	return &ErrorCode{code: code, message: msg}
}
