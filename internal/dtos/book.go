package dtos

type BookResponse struct {
	//ID              int    `json:"ID,omitempty"`
	ISBN            string `json:"ISBN"`
	UUID            string `json:"uuid,omitempty"`
	Title           string `json:"title"`
	PublicationYear int    `json:"publicationYear"`
	Description     string `json:"description,omitempty"`
	//CoverImage      string           `json:"coverImage,omitempty"`
	FilePath string              `json:"filePath"`
	Authors  BookAuthorsResponse `json:"authors"`
}

type BooksRequest struct {
	Limit int `query:"limit"`
	Page  int `query:"page"`
}

func (r *BooksRequest) Validate() error {
	r.validateLimit()
	r.validatePage()
	return nil
}

func (r *BooksRequest) validateLimit() {
	if r.Limit < 0 {
		r.Limit = 12
		return
	}

	if r.Limit == 0 {
		r.Limit = 12
		return
	}

	if r.Limit > 80 {
		r.Limit = 12
		return
	}
}

func (r *BooksRequest) validatePage() {
	if r.Page < 0 {
		r.Page = 1
		return
	}

	if r.Page == 0 {
		r.Page = 1
		return
	}

	//if r.Page > 80 {
	//	*vErr = append(*vErr, errors.New("limit must be less than 80; "))
	//	return
	//}
}

type BooksResponse struct {
	Books []BookResponse
}

type BookFileRequest struct {
	FileName string
	FileType string
	Chapter  string
}

type BookFileResponse struct {
	File []byte `json:"file"`
}
