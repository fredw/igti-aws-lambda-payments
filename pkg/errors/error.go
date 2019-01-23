package errors

// NewCriticalError returns an new critical error
func NewCriticalError(s string) error {
	return &CriticalError{s}
}

// CriticalError is an error that doesn't allow the message to be processed again
type CriticalError struct {
	s string
}

func (e *CriticalError) Error() string {
	return e.s
}
