package errors

type publicError struct {
	err error
	msg string
}

func Public(err error, msg string) error {
	return publicError{
		err: err,
		msg: msg,
	}
}

func (pe publicError) Error() string {
	return pe.err.Error()
}

func (pe publicError) Public() string {
	return pe.msg
}

func (pe publicError) Unwrap() error {
	return pe.err
}
