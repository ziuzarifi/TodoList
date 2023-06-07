package utils

import (
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	signingKey = []byte("secret")
)

func HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(hash)
}

func ComparePasswords(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func GenerateAccessToken(userID int, emailAddress string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"email_address": emailAddress,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
