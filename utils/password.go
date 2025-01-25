package utils

import "golang.org/x/crypto/bcrypt"

func GeneratePasswordHash(password string) (string , error) {
	hashedPassword , err := bcrypt.GenerateFromPassword([]byte(password) , bcrypt.DefaultCost)
	if err != nil {
		return "" , err
	}
	return string(hashedPassword) , nil
	
}

func ComparePasswordHash(password , hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword) , []byte(password))
}