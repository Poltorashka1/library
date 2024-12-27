package apperrors

import (
	"errors"
	"fmt"
)

var ErrTooShort = errors.New("parameter must be at least 1 character long")
var ErrFormatLength = errors.New("parameter must be at least 1 character long")
var ErrBookNotFound = errors.New("book not found")
var ErrBookNotExist = errors.New("book not exist")

// db errors
var ErrPageNotFound = errors.New("page not found")

type queryKeyErr struct {
	queryKey string
}

func (e *queryKeyErr) Error() string {
	return fmt.Sprintf("query key '%s' len < 1", e.queryKey)
}

func QueryKeyErr(queryKey string) error {
	return &queryKeyErr{
		queryKey,
	}
}
