package models

import (
	"errors"
	"example/buddyseller-api/database"
	"example/buddyseller-api/dtos"
	"example/buddyseller-api/utils"
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

	stmt, err := database.DB.Prepare(query)

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

func GetAllUsers() ([]User, error) {
	query := "SELECT * FROM users"
	rows, err := database.DB.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Code, &user.Email)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func GetUserById(id int64) (*User, error) {
	query := "SELECT * FROM users WHERE id = $1"
	row := database.DB.QueryRow(query, id)

	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Code, &user.Email)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (user *User) Update() error {
	query := `
	UPDATE users
	SET 
		name = $2,
		code = $3,
		email = $4
	WHERE id = $1
	`

	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(&user.ID, &user.Name, &user.Code, &user.Email)

	if err != nil {
		return err
	}

	return nil
}

func (user *User) UpdatePassword() error {
	query := `
	UPDATE users
	SET 
		password = $2,
	WHERE id = $1
	`

	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(&user.ID, &user.Password)

	if err != nil {
		return err
	}

	return nil
}

func DeleteUser(id int64) error {
	query := `
	DELETE FROM users
	WHERE id = $1
	`

	_, err := database.DB.Exec(query, id)

	if err != nil {
		return err
	}

	return nil
}

func ValidateUserCredentials(userDto *dtos.UserLoginDto) (User, error) {
	query := `
	SELECT id, email, password FROM users
	WHERE code = $1
	`

	var user User
	row := database.DB.QueryRow(query, userDto.Code)

	err := row.Scan(&user.ID, &user.Email, &user.Password)

	if err != nil {
		return User{}, err
	}

	passwordIsValid := utils.CheckPassword(userDto.Password, user.Password)

	if !passwordIsValid {
		return User{}, errors.New("credentials invalid")
	}

	user.Code = userDto.Code
	return user, nil
}
