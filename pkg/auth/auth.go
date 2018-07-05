package auth

import "golang.org/x/crypto/bcrypt"

// Encrypt encrypts the plain text with bcrypt.
func Encrypt(source string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(source), bcrypt.DefaultCost)
	return string(hashedBytes), err
}

// Compare compares the encrypted text with the plain text if it's the same.
func Compare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
