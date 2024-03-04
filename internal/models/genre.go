package models

type Genre struct {
	ID   int32  `db:"id"`
	Name string `db:"name"`
}
