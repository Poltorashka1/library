package bookusecase

import (
	"book/internal/dtos"
	"book/internal/entities"
	"context"
)

// todo add error if book file not found

func (u *bookUseCase) Book(ctx context.Context, uuid string) (*dtos.BookResponse, error) {
	result, err := u.repo.Book(ctx, uuid)
	if err != nil {
		return nil, err
	}
	// todo add converter
	// get authors from books and convert to dto
	var authors dtos.BookAuthorsResponse
	for _, author := range result.BookAuthors.Authors {
		authors.Authors = append(authors.Authors, dtos.BookAuthorResponse{
			UUID:     author.UUID,
			NickName: author.NickName,
		})
	}

	return &dtos.BookResponse{
		ISBN: result.ISBN,
		// todo mb uuid delete from response
		UUID:            result.UUID,
		Title:           result.Title,
		PublicationYear: result.PublicationYear,
		Description:     result.Description,
		FilePath:        result.BooksFileUUID.String,
		Authors:         authors,
	}, nil
}

//var author = dtos.BookAuthorResponse{}
//var authorList []dtos.BookAuthorResponse
//for _, a := range bookAuthors {
//	author.Name = a.Name
//	author.Surname = a.Surname
//	author.Patronymic = a.Patronymic
//	authorList = append(authorList, author)
//}
//
//return &dtos.BookResponse{
//	UUID:            book.UUID,
//	ISBN:            book.ISBN,
//	Title:           strings.ToTitle(book.Title), // todo check this
//	PublicationYear: book.PublicationYear,
//	Description:     book.Description,
//	FilePath:        book.BooksFileUUID,
//	Authors:         authorList,
//}, nil

func (u *bookUseCase) Books(ctx context.Context, payload *dtos.BooksRequest) (*dtos.BooksResponse, error) {
	filter := entities.BookFilter{
		Start: (payload.Page - 1) * payload.Limit,
		Stop:  (payload.Page-1)*payload.Limit + payload.Limit,
	}

	// get books list with authors
	result, err := u.repo.Books(ctx, filter)
	if err != nil {
		return nil, err
	}

	// todo add converter
	// convert books to dto
	var response dtos.BooksResponse
	for _, book := range result.Books {
		// get authors from entities.books and convert to dto
		var authors dtos.BookAuthorsResponse
		for _, author := range book.BookAuthors.Authors {
			authors.Authors = append(authors.Authors, dtos.BookAuthorResponse{
				UUID:     author.UUID,
				NickName: author.NickName,
			})
		}

		// todo сделать что то с file path
		// convert book from entities.Books and convert to dto
		//if book.BooksFileUUID == nil {
		//	*book.BooksFileUUID = ""
		//}
		response.Books = append(response.Books, dtos.BookResponse{
			UUID:     book.UUID,
			Title:    book.Title,
			FilePath: book.BooksFileUUID.String,
			Authors:  authors,
		})
	}

	return &response, nil
}
