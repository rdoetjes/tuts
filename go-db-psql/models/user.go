package models

import "phonax.com/db/rest"

type User struct {
	ID        int    `db:"id"`
	FirstName string `db:"firstname"`
	LastName  string `db:"lastname"`
	DOB       string `db:"dob"`
	Email     string `db:"email"`
}

func (u User) SetupQueries() rest.Queries {
	return rest.Queries{
		Create: `INSERT INTO users (firstname, lastname, dob, email) VALUES ($1, $2, $3, $4) RETURNING id`,
		GetOne: `SELECT id, firstname, lastname, dob, email FROM users WHERE id = $1`,
		List:   `SELECT id, firstname, lastname, dob, email FROM users`,
		Update: `UPDATE users SET firstname=$1, lastname=$2, dob=$3, email=$4 WHERE id = $5`,
		Delete: `DELETE FROM users WHERE id = $1`,
	}
}
