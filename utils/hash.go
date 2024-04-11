package utils

import "golang.org/x/crypto/bcrypt"

const complexity = 14

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), complexity)

	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
