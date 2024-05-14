package hash

import "golang.org/x/crypto/bcrypt"

func GenerateHashPassword(password string) (string, error) {
	if password == "" {
		return password, nil
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func ComparePasswords(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}
