package handlertests

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignIn(t *testing.T) {

	err := RefreshTables()
	if err != nil {
		t.Errorf("Error when refreshing tables: %v\n", err)
		return
	}

	_, _, _, _, _, _, err = SeedTables()
	if err != nil {
		log.Fatal(err)
	}

	token, err := s.SignIn("raf@bot.com", "admin")

	if err != nil {
		t.Errorf("A JWT should have been generated without error. Error raised: %v\n", err)
	}

	if token == "" {
		t.Errorf("A JWT should have been generated.")
	}

	users := []struct {
		email        string
		password     string
	}{
		{
			email:        "coach1@bot.com",
			password:     "hellogophers",
		},
		{
			email:        "client1@bot.com",
			password:     "hellogophers",
		},
		{
			email:        "coach2@bot.com",
			password:     "hellogophers",
		},
		{
			email:        "client2@bot.com",
			password:     "hellogophers",
		},
	}

	for _, v := range users {

		token, err := s.SignIn(v.email, v.password)

		if err != nil {
			t.Errorf("A JWT should have been generated without error. Error raised: %v\n", err)
		}

		if token == "" {
			t.Errorf("A JWT should have been generated.")
		}
	}
}

func TestLogin(t *testing.T) {

	err := RefreshTables()
	if err != nil {
		t.Errorf("Error when refreshing tables: %v\n", err)
		return
	}

	_, _, _, _, _, _, err = SeedTables()
	if err != nil {
		log.Fatal(err)
	}

	usersToLog := []struct {
		inputLogin   string
	}{
		{
			inputLogin:    `{"email": "raf@bot.com", "password": "admin"}`,
		},
		{
			inputLogin:    `{"email": "coach1@bot.com", "password": "hellogophers"}`,
		},
		{
			inputLogin:    `{"email": "client1@bot.com", "password": "hellogophers"}`,
		},
		{
			inputLogin:    `{"email": "coach2@bot.com", "password": "hellogophers"}`,
		},
		{
			inputLogin:    `{"email": "client2@bot.com", "password": "hellogophers"}`,
		},
	}

	for _, v := range usersToLog {

		req, err := http.NewRequest("POST", "/login", bytes.NewBufferString(v.inputLogin))
		if err != nil {
			t.Errorf("Error when sending POST request: %v", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(s.Login)
		handler.ServeHTTP(rr, req)

		if rr.Code != 200 {
			t.Errorf("Request should be completed with a 200 status code. Status code registered: %v", rr.Code)
		}
	}
}
