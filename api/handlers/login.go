package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Alexandremerancienne/my_Sartorius/api/auth"
	"github.com/Alexandremerancienne/my_Sartorius/api/exceptions"
	"github.com/Alexandremerancienne/my_Sartorius/api/models"
	"golang.org/x/crypto/bcrypt"
)

// Login logs user after successful authentication.
func (server *Server) Login(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		exceptions.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		err = errors.New("422 Unprocessable Entity Error: Cannot parse JSON data.")
		exceptions.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	token, err := server.SignIn(user.Email, user.Password)
	if err != nil {
		err = errors.New("Login has failed: please try again.")
		exceptions.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(token)
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}

// SignIn generates a JSON Web Token after successful authentication.
func (server *Server) SignIn(email, password string) (string, error) {

	var err error
	user := models.User{}

	err = server.Database.Debug().Model(models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return "", err
	}
	err = models.CheckPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	return auth.GenerateJWT(user.ID)
}
