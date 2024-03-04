package models

type Country struct {
	ID   int32  `db:"id"`
	Name string `db:"name"`
}
