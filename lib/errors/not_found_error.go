package errors

import "errors"

type NotFoundError struct {
	message string
}

func (e NotFoundError) Error() string {
	return e.message
}

// NotFoundErrorを生成するヘルパー関数
func NewNotFoundError(message string) error {
	return NotFoundError{message: message}
}

func IsNotFoundError(err error) bool {
	var notFoundError NotFoundError
	return errors.As(err, &notFoundError)
}
