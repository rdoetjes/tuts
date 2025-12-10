package models

import "phonax.com/db/rest"

type User struct {
	ID        int    `db:"id"`
	FirstName string `db:"firstname"`
	LastName  string `db:"lastname"`
	DOB       string `db:"dob"`
	Email     string `db:"email"`
	Password  string `db:"password"`
}

func (u User) SetupQueries() rest.Queries {
	return rest.Queries{
		Create: `INSERT INTO users (firstname, lastname, dob, email, password) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		GetOne: `SELECT id, firstname, lastname, dob, email FROM users WHERE id = $1`,
		List:   `SELECT id, firstname, lastname, dob, email FROM users`,
		Update: `UPDATE users SET firstname=$1, lastname=$2, dob=$3, email=$4, password=$5 WHERE id = $6`,
		Delete: `DELETE FROM users WHERE id = $1`,
	}
}
