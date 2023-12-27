package errs

type ValidateError struct {
	errorMsg string
}

func NewValidateError(errorMsg string) *ValidateError {
	return &ValidateError{errorMsg: errorMsg}
}

func (v *ValidateError) Error() string {
	return v.errorMsg
}
