package request

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	// ErrUnknownContentType returns if request content type is not supported
	ErrUnknownContentType = errors.New("unknown content type")
	ErrInvalidJsonSyntax  = errors.New("invalid JSON syntax")

	// ErrContentToLarge returns if content size more then allowed
	//ErrContentToLarge = errors.New("content size to large")

	// ErrFieldLength returns if file size more then allowed
	ErrFieldLength      = errors.New("file length exceeds the limit")
	ErrInvalidFieldType = errors.New("only 'file' field can contain file content")
	ErrFileNameTooLong  = errors.New("file name too long, max length 100")
	//ErrFieldName          = errors.New("invalid field name")
)

// ErrContentToLarge returns if content size more than allowed in mb
type ErrContentToLarge struct {
	limit int
}

func (e *ErrContentToLarge) Error() string {
	return fmt.Sprintf("content too large, max size %d mb", e.limit/(1024*1024))
}

// ErrFileToLarge returns if file size more than allowed in mb
type ErrFileToLarge struct {
	limit int
}

func (e *ErrFileToLarge) Error() string {
	return fmt.Sprintf("file size too large, max size %d mb", e.limit/(1024*1024))
}

type ErrInvalidFileType struct {
	fileType string
}

func (e *ErrInvalidFileType) Error() string {
	return fmt.Sprintf("invalid file type, only %s; ", e.fileType)
}

// ErrFieldName returns if form field name is missing from the source struct tagNames
type ErrFieldName struct {
	fieldName string
}

func (e *ErrFieldName) Error() string {
	return fmt.Sprintf("invalid field name '%s'", e.fieldName)
}

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

func (e *MultiError) Add(err error) {
	e.err = append(e.err, err)
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
