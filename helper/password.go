package helper

import (
	"golang.org/x/crypto/bcrypt"
)

// EncryptPassword encrypts the password given to the function
func EncryptPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

// DecryptPassword compares the hashedPassword to the password and see if there is a match
func DecryptPassword(password string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword) , []byte(password))
	if err != nil {
		return false
	}

	return true
}