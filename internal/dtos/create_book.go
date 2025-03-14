package dtos

import (
	apperrors "book/internal/errors"
	"errors"
	"fmt"
	"os"
	"strings"
	"unicode"
)

// todo обезопасить то большого количества данных в строках
// todo разные форматы isbn
// todo test form parsing с использованием встраивания

type CreateBookRequest struct {
	BookInfo `json:"book_info"`
	Files    Files `json:"files"`
	File     File  `json:"file"`
	//Digits   []int   `json:"digits" form:"digits"`
	//Authors  Authors `json:"authors" form:"authors"`
}

//type Authors []Author

type Author struct {
	NickName string `json:"nick_name" form:"nick_name,required"`
	Name     string `json:"name" form:"name,optional"`
	Surname  string `json:"surname" form:"surname,optional"`
}

type Files struct {
	Files []*os.File `form:"files"`
}
type File struct {
	File *os.File `form:"file,required"`
}

type BookInfo struct {
	Title           string `json:"title" form:"title,required"`
	ISBN            string `form:"isbn,required"`
	PublicationYear int    `json:"publication_year,required" form:"publication_year,optional"`
	Description     string `json:"description" form:"description,"`
	Publisher       string `json:"publisher" form:"publisher,required"`
}

type CreateBookResponse struct {
	ISBN string `json:"isbn"`
}

func (r *CreateBookRequest) Validate() error {
	var e apperrors.ValidationErrors

	err := r.titleValidate()
	if err != nil {
		e = append(e, err)
	}
	err = r.isbnValidate()
	if err != nil {
		e = append(e, err)
	}
	err = r.publicationYearValidate()
	if err != nil {
		e = append(e, err)
	}
	//err = r.descriptionValidate()
	//if err != nil {
	//	e = append(e, err)
	//}
	//err = r.publisherValidate()
	//if err != nil {
	//	e = append(e, err)
	//}
	if e != nil {
		return e
	}
	return nil
}

func (r *CreateBookRequest) titleValidate() error {
	if r.Title == "" {
		return errors.New("empty field Title; ")
	}

	if len(r.Title) > 255 {
		return errors.New("title too long; ")
	}
	return nil
}

func (r *CreateBookRequest) isbnValidate() error {
	r.ISBN = strings.ReplaceAll(r.ISBN, "-", "")
	r.ISBN = strings.ReplaceAll(r.ISBN, " ", "")

	if r.ISBN == "" {
		return errors.New("ISBN is empty")
	}

	if len(r.ISBN) == 10 {
		if !r.ISValidISBN10() {
			return errors.New("invalid ISBN; ")
		}

		r.ISBN10To13()
	}

	if len(r.ISBN) != 13 {
		return errors.New("invalid ISBN; ")
	}

	if !r.ISValidISBN13() {
		return errors.New("invalid ISBN; ")
	}

	return nil
}

func (r *CreateBookRequest) ISBN10To13() {
	isbn13 := "978" + r.ISBN[:9]

	sum := 0
	for i, r := range isbn13 {
		digit := int(r - '0')
		if i%2 == 0 {
			sum += digit
		} else {
			sum += digit * 3
		}
	}
	checkDigit := (10 - (sum % 10)) % 10

	r.ISBN = isbn13 + fmt.Sprintf("%d", checkDigit)
}

func (r *CreateBookRequest) ISValidISBN13() bool {
	sum := 0
	for i, r := range r.ISBN {
		if !unicode.IsDigit(r) {
			return false
		} else {
			digit := int(r - '0')
			if i%2 == 0 {
				sum += digit
			} else {
				sum += digit * 3
			}
		}
	}

	return sum%10 == 0
}

func (r *CreateBookRequest) ISValidISBN10() bool {
	sum := 0
	for i, char := range r.ISBN {
		var digit int
		if i == 9 && (char == 'X' || char == 'x') {
			digit = 10
		} else if !unicode.IsDigit(char) {
			return false
		} else {
			digit = int(char - '0')
		}
		sum += digit * (10 - i)
	}

	return sum%11 == 0
}

func (r *CreateBookRequest) publicationYearValidate() error {
	if r.BookInfo.PublicationYear <= 0 {
		return errors.New("publication year must be greater than 0; ")
	}
	return nil
}

//func (r *CreateBookRequest) descriptionValidate() error {
//	return nil
//}

//func (r *CreateBookRequest) publisherValidate() error {
//	return nil
//}
