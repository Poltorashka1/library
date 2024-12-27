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
	Limit int
	Page  int
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
