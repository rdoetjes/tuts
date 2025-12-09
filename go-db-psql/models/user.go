package models

type User struct {
	ID        int    `db:"id"`
	FirstName string `db:"firstname"`
	LastName  string `db:"lastname"`
	DOB       string `db:"dob"`
	Email     string `db:"email"`
}
