package entities

import "database/sql"

// todo delete json tags

type Author struct {
	ID         int
	UUID       string
	NickName   string
	Name       sql.NullString
	Surname    sql.NullString
	Patronymic sql.NullString
}

type Authors []Author

//type Authors struct {
//	Authors []Author
//}

//type Authors struct {
//	Authors []Author `json:"authors"`
//}
