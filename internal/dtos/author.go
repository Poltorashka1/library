package dtos

type BookAuthorResponse struct {
	UUID       string `json:"uuid"`
	Name       string `json:"name,omitempty"`
	Surname    string `json:"surname,omitempty"`
	Patronymic string `json:"patronymic,omitempty"`
	NickName   string `json:"nickName"`
}

type BookAuthorsResponse struct {
	Authors []BookAuthorResponse `json:"authors"`
}

//	type BooksAuthorsResponse struct {
//		Authors []BooksAuthorResponse `json:"authors"`
//	}

//type BooksAuthorResponse struct {
//	UUID     string `json:"uuid"`
//	NickName string `json:"NickName"`
//}
