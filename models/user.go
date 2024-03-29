package models

import (
	"example/buddyseller-api/db"
)

type User struct {
	ID       int64
	Name     string `binding:"required"`
	Code     string `binding:"required"`
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func (user *User) Save() error {
	query := `
	INSERT INTO users(name, code, email, password)
	VALUES ($1, $2, $3, $4) RETURNING id
	`

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()
	stmt.Exec()
	var lastInsertedId int64
	err = stmt.QueryRow(&user.Name, &user.Code, &user.Email, &user.Password).Scan(&lastInsertedId)

	if err != nil {
		return err
	}

	user.ID = lastInsertedId
	return err
}
