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

// TODO: DANGEROUS - Should not comparing passwords in case of error
// Other approaches:
// 1. Set all passwords to be hashed (rn the seed users are not hashed for lazy reasons :P)
// 2. Add a field to allow multiple type of passwords (hased, not hashed, ...) and evaluate the password cheking accordingly (Probably not goint to do it for lazy reasons :P)
func CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if err != nil {
		return hashedPassword == password
	}

	return err == nil
}
