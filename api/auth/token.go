package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/joho/godotenv"
)

var SecretKey = []byte{}

// returnKey returns the value of SECRET_KEY environment variable.
func returnKey() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Error loading .env")
	}
	key := os.Getenv("SECRET_KEY")
	SecretKey = []byte(key)
}

// GenerateJWT generates and configures JSON Web Token (JWT).
func GenerateJWT(user_id uint) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user_id
	claims["exp"] = time.Now().Add(10 * time.Minute).Unix()
	claims["authorized"] = true
	returnKey()
	tokenString, err := token.SignedString(SecretKey)
	return tokenString, err
}

// CheckJWT checks is the JWT provided is valid.
func CheckJWT(r *http.Request) bool {

	tokenString := ExtractToken(r)
	extractedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Wrong Signing Method")
		}
		return SecretKey, nil
	})
	if err != nil {
		return false
	}
	return extractedToken.Valid
}

// ExtractToken extracts the JWT from a request.
func ExtractToken(r *http.Request) string {
	keys := r.URL.Query()
	token := keys.Get("token")
	if token != "" {
		return token
	}
	bearerToken := r.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

// ExtractTokenID extracts thd ID of the user using a JWT.
func ExtractTokenID(r *http.Request) (uint32, error) {

	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["user_id"]), 10, 32)
		if err != nil {
			return 0, err
		}
		return uint32(uid), nil
	}
	return 0, nil
}
