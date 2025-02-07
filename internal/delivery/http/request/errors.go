package request

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrUnknownContentType = errors.New("unknown content type")
	ErrInvalidJsonSyntax  = errors.New("invalid JSON syntax")
	ErrInvalidContentType = errors.New("invalid file type, only PDF and EPUB; ") // todo add zip or idn
	ErrContentToLarge     = errors.New("content size to large")
	ErrFieldLength        = errors.New("file length exceeds the limit")
	ErrInvalidFieldType   = errors.New("only 'file' field can contain file content")
	ErrFileNameTooLong    = errors.New("file name too long, max length 100")
	ErrFieldName          = errors.New("invalid field name")
)

type ErrFieldType struct {
	field string
	typ   reflect.Type
}

func (e *ErrFieldType) Error() string {
	return fmt.Sprintf("the field '%s' must be of type %s", e.field, e.typ)
}

type ErrFieldRequired struct {
	field string
}

func (e ErrFieldRequired) Error() string {
	return fmt.Sprintf("the field '%s' is required; ", e.field)
}

type MultiError struct {
	err []error
}

func (e *MultiError) Error() string {
	var err string
	for _, e := range e.err {
		err += e.Error()
	}
	return err
}

type ErrQueryRequired struct {
	queryKey string
}

func (e ErrQueryRequired) Error() string {
	return fmt.Sprintf("the query '%s' is required; ", e.queryKey)
}
