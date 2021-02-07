package models

type Contact struct {
	Name string
	UserID int
	Telephone string `db:"tel"`
	Description string
}