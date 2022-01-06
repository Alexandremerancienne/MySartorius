package handlertests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Alexandremerancienne/my_Sartorius/api/models"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUsers_WhenGetUsersAsAdmin_AllUsersAreAvailable(t *testing.T) {

	initializeTestData(t)

	token, err := s.SignIn("raf@bot.com", "admin")
	if err != nil {
		log.Fatalf("Error when login: %v\n", err)
	}

	req, err := http.NewRequest("GET", "/users", nil)
	req.Header.Set("Authorization", "Bearer " + token)
	if err != nil {
		t.Errorf("Error when getting users: %v\n", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.GetUsers)
	handler.ServeHTTP(rr, req)

	var users []models.User

	err = json.Unmarshal([]byte(rr.Body.String()), &users)
	if err != nil {
		log.Fatalf("Cannot convert to JSON: %v\n", err)
	}

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(users), 5)
}

func TestGetUser_WhenGetUserAsAdmin_UserIsAvailable(t *testing.T) {

	initializeTestData(t)

	token, err := s.SignIn("raf@bot.com", "admin")
	if err != nil {
		log.Fatalf("Error when login: %v\n", err)
	}

	users := []struct {
		id 	string
		first_name string
		last_name string
		email string
	}{
		{
			id: "1",
			first_name: "Raphaël",
			last_name: "Merancienne",
			email: "raf@bot.com",
		},
		{
			id: "2",
			first_name: "Coach 1",
			last_name: "Whatever",
			email: "coach1@bot.com",
		},
		{
			id: "3",
			first_name: "Client 1",
			last_name: "Whatever",
			email: "client1@bot.com",
		},
		{
			id: "4",
			first_name: "Coach 2",
			last_name: "Whatever",
			email: "coach2@bot.com",
		},
		{
			id: "5",
			first_name: "Client 2",
			last_name: "Whatever",
			email: "client2@bot.com",
		},
	}

	for _, user := range users {

		req, err := http.NewRequest("GET", "/users/" + user.id, nil)
		req.Header.Set("Authorization", "Bearer " + token)
		if err != nil {
			t.Errorf("Error when getting user: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": user.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(s.GetUser)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to JSON: %v\n", err)
		}

		assert.Equal(t, rr.Code, 200)
		assert.Equal(t, user.first_name, responseMap["first_name"])
		assert.Equal(t, user.last_name, responseMap["last_name"])
		assert.Equal(t, user.email, responseMap["email"])
	}
}

func TestGetUsers_WhenGetUsersAsCoachOrClient_OnlyRequestUserIsAvailable(t *testing.T) {

	initializeTestData(t)

	users := []struct {
		email string
		userIdRequired string
	}{
		{
			email: "coach1@bot.com",
			userIdRequired: "3",
		},
		{
			email: "client1@bot.com",
			userIdRequired: "4",
		},
		{
			email: "coach2@bot.com",
			userIdRequired: "5",
		},
		{
			email: "client2@bot.com",
			userIdRequired: "4",
		},
	}

	for _, user := range users {

		token, err := s.SignIn(user.email, "hellogophers")
		if err != nil {
			log.Fatalf("Error when login: %v\n", err)
		}

		req, err := http.NewRequest("GET", "/users/" + user.userIdRequired, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		if err != nil {
			t.Errorf("Error when getting users: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": user.userIdRequired})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(s.GetUser)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, http.StatusUnauthorized)
	}
}

func TestGetUser_WhenGetUserAsCoachOrClientOtherThanRequestUser_401ErrorIsReturned(t *testing.T) {

	initializeTestData(t)

	u := models.User{}
	allUsers, err := u.GetUsers(s.Database)

	if err != nil {
		t.Errorf("Error when getting users: %v\n", err)
		return
	}

	for _, user := range *allUsers {

		var finalUsers []models.User

		if user.Role != "manager" {
			token, err := s.SignIn(user.Email, "hellogophers")
			if err != nil {
				log.Fatalf("Error when login: %v\n", err)
			}

			req, err := http.NewRequest("GET", "/users", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			if err != nil {
				t.Errorf("Error when getting users: %v\n", err)
			}
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(s.GetUsers)
			handler.ServeHTTP(rr, req)

			var u = models.User{}
			err = json.Unmarshal([]byte(rr.Body.String()), &u)
			if err != nil {
				log.Fatalf("Cannot convert to JSON: %v\n", err)
			}

			assert.Equal(t, rr.Code, http.StatusOK)

			finalUsers = append(finalUsers, u)
		} else {
			continue
		}

		assert.Equal(t, len(finalUsers), 1)
	}
}

func TestCreateUser_WhenCreateUserAsAdmin_UserIsCreated(t *testing.T) {

	initializeTestData(t)

	token, err := s.SignIn("raf@bot.com", "admin")
	if err != nil {
		log.Fatalf("Error when login: %v\n", err)
	}

	usersToCreate := []struct {
		inputJSON    string
		statusCode   int
		firstName string
		lastName string
		email string
	}{
		{
			inputJSON:    `{"first_name": "Vladimir", "last_name": "Petrov", "email": "vladimir_client@bot.com",
							"password": "hellogophers", "phone_number": "+359452154050", "role":"client"}`,
			statusCode:   201,
			firstName: "Vladimir",
			lastName: "Petrov",
			email: "vladimir_client@bot.com",
		},
		{
			inputJSON:    `{"first_name": "Emil", "last_name": "Petrov", "email": "vladimir_client@bot.com",
							"password": "hellogophers", "phone_number": "+359452154050", "role":"client"}`,
			statusCode:   400,
		},
		{
			inputJSON:    `{"first_name": "", "last_name": "Petrov", "email": "emil_client@bot.com",
							"password": "hellogophers", "phone_number": "+359452154050", "role":"client"}`,
			statusCode:   400,
		},
		{
			inputJSON:    `{"first_name": "Tatiana", "last_name": "", "email": "tatiana_client@bot.com",
							"password": "hellogophers", "phone_number": "+359452154050", "role":"coach"}`,
			statusCode:   400,
		},
	}

	for _, user := range usersToCreate {

		req, err := http.NewRequest("POST", "/users", bytes.NewBufferString(user.inputJSON))
		if err != nil {
			t.Errorf("Error when creating user: %v", err)
		}
		req.Header.Set("Authorization", "Bearer " + token)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(s.CreateUser)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to JSON: %v", err)
		}
		assert.Equal(t, rr.Code, user.statusCode)
		if user.statusCode == 201 {
			assert.Equal(t, responseMap["first_name"], user.firstName)
			assert.Equal(t, responseMap["last_name"], user.lastName)
			assert.Equal(t, responseMap["email"], user.email)
		}
	}
}

func TestCreateUser_WhenCreateUserAsCoachOrClient_401ErrorIsReturned(t *testing.T) {

	initializeTestData(t)

	users := []struct {
		email    string
		password     string
	}{
		{
			email: "coach1@bot.com",
			password: "hellogophers",
		},
		{
			email: "coach2@bot.com",
			password: "hellogophers",
		},
		{
			email: "client1@bot.com",
			password: "hellogophers",
		},
		{
			email: "client2@bot.com",
			password: "hellogophers",
		},
	}

	for _, user := range users {

		token, err := s.SignIn(user.email, user.password)
		if err != nil {
			log.Fatalf("Error when login: %v\n", err)
		}

		req, err := http.NewRequest("POST", "users/", nil)
		if err != nil {
			t.Errorf("Error when creating user: %v\n", err)
		}
		req.Header.Set("Authorization", "Bearer "+token)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(s.CreateUser)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, http.StatusUnauthorized)
	}
}

func TestUpdateUser_WhenUpdateUserAsCoachOrClient_401ErrorIsReturned(t *testing.T) {

	initializeTestData(t)

	users := []struct {
		email    string
		password     string
		userUpdated string
	}{
		{
			email: "coach1@bot.com",
			password: "hellogophers",
			userUpdated: "1",
		},
		{
			email: "coach2@bot.com",
			password: "hellogophers",
			userUpdated: "2",
		},
		{
			email: "client1@bot.com",
			password: "hellogophers",
			userUpdated: "1",
		},
		{
			email: "client2@bot.com",
			password: "hellogophers",
			userUpdated: "2",
		},
	}

	for _, user := range users {

		token, err := s.SignIn(user.email, user.password)
		if err != nil {
			log.Fatalf("Error when login: %v\n", err)
		}

			req, err := http.NewRequest("PUT", "users/" + user.userUpdated, nil)
			if err != nil {
				t.Errorf("Error when updating user: %v\n", err)
			}
			req.Header.Set("Authorization", "Bearer "+token)
			req = mux.SetURLVars(req, map[string]string{"id": user.userUpdated})

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(s.UpdateUser)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, rr.Code, http.StatusUnauthorized)
	}
}

func TestDeleteUser_WhenDeleteUserAsAdmin_UserIsDeleted(t *testing.T) {

	initializeTestData(t)

	token, err := s.SignIn("raf@bot.com", "admin")
	if err != nil {
		log.Fatalf("Error when login: %v\n", err)
	}

	req, err := http.NewRequest("DELETE", "/users", nil)
	if err != nil {
		t.Errorf("Error when deleting users: %v\n", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req = mux.SetURLVars(req, map[string]string{"id": "2"})
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.DeleteUser)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusNoContent)

	req, err = http.NewRequest("GET", "/users", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	if err != nil {
		t.Errorf("Error when getting users: %v\n", err)
	}
	rr = httptest.NewRecorder()
	handler = s.GetUsers
	handler.ServeHTTP(rr, req)

	var users = []models.User{}
	err = json.Unmarshal([]byte(rr.Body.String()), &users)
	if err != nil {
		log.Fatalf("Cannot convert to JSON: %v\n", err)
	}

	assert.Equal(t, len(users), 4)
}

func TestDeleteUser_WhenDeleteUserAsCoachOrClient_401ErrorIsReturned(t *testing.T) {

	initializeTestData(t)

	users := []struct {
		email    string
		password     string
	}{
		{
			email: "coach1@bot.com",
			password: "hellogophers",
		},
		{
			email: "coach2@bot.com",
			password: "hellogophers",
		},
		{
			email: "client1@bot.com",
			password: "hellogophers",
		},
		{
			email: "client2@bot.com",
			password: "hellogophers",
		},
	}

	usersIDs := []string{"1", "2", "3", "4", "5"}

	for _, user := range users {

		token, err := s.SignIn(user.email, user.password)
		if err != nil {
			log.Fatalf("Error when login: %v\n", err)
		}

		for _, userID := range usersIDs {
			req, err := http.NewRequest("DELETE", "users", nil)
			if err != nil {
				t.Errorf("Error when deleting user: %v\n", err)
			}
			req.Header.Set("Authorization", "Bearer "+token)
			req = mux.SetURLVars(req, map[string]string{"id": userID})

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(s.DeleteUser)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, rr.Code, http.StatusUnauthorized)
		}
	}
}
