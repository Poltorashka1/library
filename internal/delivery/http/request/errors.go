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
	ErrFileNameTooLong    = errors.New("file name too long, max length 100")
)

//type UnprocessableEntity interface {
//	UnprocessableEntity()
//}

// ErrContentToLarge returns if content size more than allowed in bytes/kb/mb
type ErrContentToLarge struct {
	limit int
}

func (e *ErrContentToLarge) Error() string {
	// return ib bytes
	if e.limit < 1024 {
		return fmt.Sprintf("the total size of thr transmitted data exceeds the maximum limit in %d bytes", e.limit)
	}
	// return ib kb
	if e.limit < 1048576 {
		return fmt.Sprintf("the total size of thr transmitted data exceeds the maximum limit in %d kb", e.limit/1024)
	}
	// return ib mb
	return fmt.Sprintf("the total size of thr transmitted data exceeds the maximum limit in %d mb", e.limit/(1024*1024))
}

// ErrFormValueToLarge returns if form field value size more than allowed in bytes/kb/mb
type ErrFormValueToLarge struct {
	formField string
	limit     int
}

func (e *ErrFormValueToLarge) Error() string {
	// return ib bytes
	if e.limit < 1024 {
		return fmt.Sprintf("the size of the '%s' field has been exceeded, max size %d bytes", e.formField, e.limit)
	}
	// return ib kb
	if e.limit < 1048576 {
		return fmt.Sprintf("the size of the '%s' field has been exceeded, max size %d kb", e.formField, e.limit/1024)
	}
	// return ib mb
	return fmt.Sprintf("the size of the '%s' field has been exceeded, max size %d mb", e.formField, e.limit/(1024*1024))
}

// ErrInvalidFileType returns if form file content type is not supported
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

//func (e *ErrFieldName) UnprocessableEntity() {}

// ErrFieldType returns if field type
type ErrFieldType struct {
	field string
	typ   reflect.Type
}

func (e *ErrFieldType) Error() string {
	return fmt.Sprintf("the field '%s' must be of type %s", e.field, e.typ)
}

//type ErrFieldRequired struct {
//	field string
//}
//
//func (e ErrFieldRequired) Error() string {
//	return fmt.Sprintf("the field '%s' is required; ", e.field)
//}

// MultiError returns all field validation errors for example fieldType/fieldRequired
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

//func (e *MultiError) UnprocessableEntity() {}

type ErrQueryRequired struct {
	queryKey string
}

func (e ErrQueryRequired) Error() string {
	return fmt.Sprintf("the query '%s' is required; ", e.queryKey)
}
