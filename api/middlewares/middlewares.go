package middlewares

import (
	"errors"
	"github.com/Alexandremerancienne/my_Sartorius/api/auth"
	"github.com/Alexandremerancienne/my_Sartorius/api/exceptions"
	"net/http"
)

// SetMiddlewareAuthentication checks for the validity of the authentication token provided.
func SetMiddlewareAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b := auth.CheckJWT(r)
		if !b {
			err := errors.New("Invalid token: please login with valid credentials")
			exceptions.ERROR(w, http.StatusUnauthorized, err)
			return
		}
		next(w, r)
	}
}

