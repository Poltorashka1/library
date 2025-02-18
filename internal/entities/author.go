package entities

// todo delete json tags

type Author struct {
	ID         int
	UUID       string
	NickName   string
	Name       string
	Surname    string
	Patronymic *string
}

type Authors []Author

//type Authors struct {
//	Authors []Author
//}

//type Authors struct {
//	Authors []Author `json:"authors"`
//}
