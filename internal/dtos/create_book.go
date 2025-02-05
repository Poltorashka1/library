package dtos

import (
	apperrors "book/internal/errors"
	"errors"
	"os"
	"regexp"
	"strings"
)

// todo обезопасить то большого количества данных в строках
// todo разные форматы isbn
type CreateBookRequest struct {
	BookInfo BookInfo `json:"book-info" form:"book-info"`
	Files    Files    `json:"files" form:"files"`
}

type Files struct {
	File  *os.File   `form:"file,required"`
	Files []*os.File `form:"files,required"`
}

type BookInfo struct {
	Title           string `json:"title" form:"title,required"`
	ISBN            string `json:"isbn" form:"isbn,required"`
	PublicationYear int    `json:"publication-year,required" form:"publication-year,required"`
	Description     string `json:"description" form:"description,optional"`
	Publisher       string `json:"publisher" form:"publisher,required"`
}

type CreateBookResponse struct {
	ISBN string `json:"isbn"`
}

func (r *CreateBookRequest) Validate() error {
	e := &apperrors.ValidationErrors{}

	//r.titleValidate(e)
	r.isbnValidate(e)
	r.publicationYearValidate(e)
	//r.descriptionValidate(e)
	//r.publisherValidate(e)

	if e.Errors != nil {
		return e
	}
	return nil
}

//func (r *CreateBookRequest) titleValidate(e *apperrors.ValidationErrors) {
//	if r.Title == "" {
//		e.Errors = append(e.Errors, errors.New("empty field Title; "))
//		return
//	}
//
//	r.Title = optimizeString(r.Title)
//
//	if len(r.Title) > 255 {
//		e.Errors = append(e.Errors, errors.New("title too long; "))
//		return
//	}
//
//}

func optimizeString(s string) string {
	newStr := strings.TrimSpace(s)
	re := regexp.MustCompile(`\s{2,}`)
	newStr = re.ReplaceAllString(newStr, " ")
	return newStr
}

func (r *CreateBookRequest) isbnValidate(e *apperrors.ValidationErrors) {
	if r.BookInfo.ISBN == "" {
		e.Errors = append(e.Errors, errors.New("empty field ISBN; "))
		return
	}

	r.BookInfo.ISBN = optimizeString(r.BookInfo.ISBN)

}

func (r *CreateBookRequest) publicationYearValidate(e *apperrors.ValidationErrors) {
	if r.BookInfo.PublicationYear <= 0 {
		e.Errors = append(e.Errors, errors.New("publication year must be greater than 0; "))
		return
	}

}
