package errors

import (
	"fmt"
	"io"
)

// WrapOrNew returns a wrapped error if err exists or a new error with the supplied message.
func WrapOrNew(err error, message string) error {
	if err == nil {
		return New(message)
	}

	return Wrap(err, message)
}

// WrapOrErrorf formats according to a format specifier and returns the string
// as a value that satisfies error or returns an error annotating err with a stack trace
// at the point WrapOrErrorf is called, and the format specifier.
func WrapOrErrorf(err error, format string, args ...interface{}) error {
	if err == nil {
		return Errorf(format, args...)
	}

	return Wrapf(err, format, args...)
}

type withFields struct {
	fields []interface{}
	error
}

// WithFields annotates err with the given fields at the point WithFields was called.
func WithFields(err error, keyValues ...interface{}) error {
	if fErr, ok := err.(*withFields); ok {
		return &withFields{
			fields: append(fErr.fields, keyValues...),
			error:  fErr.error,
		}
	}

	return &withFields{
		fields: keyValues,
		error:  err,
	}
}

func (w *withFields) Fields() []interface{} {
	return w.fields
}

func (w *withFields) Cause() error {
	return w.error
}

func (w *withFields) UnWrap() error {
	return w.error
}

func (w *withFields) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "%+v\n", w.Cause())
			_, _ = io.WriteString(s, w.Error())
			return
		}
		fallthrough
	case 's', 'q':
		_, _ = io.WriteString(s, w.Error())
	}
}

func Field(err error, key interface{}) interface{} {
	type causer interface {
		Cause() error
	}

	type fielder interface {
		Fields() []interface{}
	}

	for err != nil {
		if fErr, ok := err.(fielder); ok {
			fields := fErr.Fields()
			for i := 0; i < len(fields)-1; i += 2 {
				if fields[i] == key {
					return fields[i+1]
				}
			}
		}

		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}

	return nil
}
