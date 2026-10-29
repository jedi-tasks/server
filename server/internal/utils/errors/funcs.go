package errors

import "pkg/errors"

func Is(err error, target error) bool {
	return errors.Is(err, target)
}

func As(get error, target any) bool {
	return errors.As(get, target)
}

func New(text string) error {
	return errors.New(text)
}

func NewMultiError() *errors.MultiError {
	return errors.NewMultiError()
}

func CastError(err error) errors.Error {
	return errors.CastError(err)
}

func IsContextError(err error) bool {

	var customErr errors.Error
	if As(err, &customErr) {
		if customErr.ErrorType == ContextCancelled.ErrorType {
			return true
		}
	}

	return errors.IsContextError(err)
}

func IsInternal(err error) bool {
	var customErr errors.Error
	if As(err, &customErr) {
		return customErr.ErrorType == InternalServer.ErrorType
	}
	return true
}

func HasInternal(errs []error) bool {
	for _, err := range errs {
		if IsInternal(err) {
			return true
		}
	}
	return false
}

func HasContextError(errs []error) bool {
	for _, err := range errs {
		if IsContextError(err) {
			return true
		}
	}
	return false
}

const (
	LogAsError   = errors.LogAsError
	LogAsWarning = errors.LogAsWarning
	LogAsDebug   = errors.LogAsDebug
	LogAsInfo    = errors.LogAsInfo
	LogNone      = errors.LogNone
)
