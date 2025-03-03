package apperrors

import (
	"errors"
	"fmt"
)

var ErrTooShort = errors.New("parameter must be at least 1 character long")
var ErrFormatLength = errors.New("parameter must be at least 1 character long")
var ErrBookNotFound = errors.New("book not found")
var ErrBookNotExist = errors.New("book not exist")

var ErrServerError = errors.New("server error")

// debug

//var (
//	NonPointerError = errors.New("transmitted value need to be pointer")
//	NilPointerError = errors.New("transmitted pointer indicate to nil")
//)

// db errors
var ErrPageNotFound = errors.New("page not found")

type ErrContentToLarge struct {
	MaxSize int64
}

func (e *ErrContentToLarge) Error() string {
	return fmt.Sprintf("content too large, max size %d mb", e.MaxSize/(1024*1024))
}

type ValidationErrors struct {
	Errors []error
}

//type ParseErrors struct {
//	Errors []error
//}
//
//func (e *ParseErrors) Error() string {
//	var err string
//
//	for _, e := range e.Errors {
//		err += e.Error()
//	}
//
//	return err
//}

// todo стоит оптимизировать потому что создается много объектов, кроме того оптимизировать обновление строки
func (e *ValidationErrors) Error() string {
	var err string

	for _, e := range e.Errors {
		err += e.Error()
	}

	return err
}

//type validationError struct {
//	error string
//}
//
//func (e *validationError) Error() string {
//	return e.error
//}
//
//func ValidationError(err string) error {
//	return &validationError{
//		err,
//	}
//}

//type queryKeyErr struct {
//	queryKey string
//}
//
//func (e *queryKeyErr) Error() string {
//	return fmt.Sprintf("query key '%s' len < 1", e.queryKey)
//}
//
//func QueryKeyErr(queryKey string) error {
//	return &queryKeyErr{
//		queryKey,
//	}
//}
