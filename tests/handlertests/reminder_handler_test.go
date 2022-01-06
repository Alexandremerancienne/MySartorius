package handlertests

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetReminders_WhenGetRemindersAsAdminOrCoach_401ErrorIsReturned(t *testing.T) {

	initializeTestData(t)

	users := []struct {
		email    string
		password     string
	}{
		{
			email: "raf@bot.com",
			password: "admin",
		},
		{
			email: "coach1@bot.com",
			password: "hellogophers",
		},
		{
			email: "coach2@bot.com",
			password: "hellogophers",
		},
	}

	for _, user := range users {
		token, err := s.SignIn(user.email, user.password)
		if err != nil {
			log.Fatalf("Error when login: %v\n", err)
		}

		req, err := http.NewRequest("GET", "tasks/1/reminders", nil)
		if err != nil {
			t.Errorf("Error when getting reminders: %v\n", err)
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req = mux.SetURLVars(req, map[string]string{"task_id": "1"})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(s.GetReminders)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, http.StatusUnauthorized)
	}
}

func TestGetReminder_WhenGetReminderAsAdminOrCoach_401ErrorIsReturned(t *testing.T) {

	initializeTestData(t)

	users := []struct {
		email    string
		password     string
	}{
		{
			email: "raf@bot.com",
			password: "admin",
		},
		{
			email: "coach1@bot.com",
			password: "hellogophers",
		},
		{
			email: "coach2@bot.com",
			password: "hellogophers",
		},
	}

	for _, user := range users {
		token, err := s.SignIn(user.email, user.password)
		if err != nil {
			log.Fatalf("Error when login: %v\n", err)
		}

		req, err := http.NewRequest("GET", "tasks/1/reminders/1", nil)
		if err != nil {
			t.Errorf("Error when getting reminders: %v\n", err)
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req = mux.SetURLVars(req, map[string]string{"task_id": "1", "reminder_id": "1"})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(s.GetReminders)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, http.StatusUnauthorized)
	}
}

func TestUpdateReminder_WhenUpdateReminderAsAdminOrCoach_401ErrorIsReturned(t *testing.T) {

	initializeTestData(t)

	users := []struct {
		email    string
		password     string
	}{
		{
			email: "raf@bot.com",
			password: "admin",
		},
		{
			email: "coach1@bot.com",
			password: "hellogophers",
		},
		{
			email: "coach2@bot.com",
			password: "hellogophers",
		},
	}

	for _, user := range users {
		token, err := s.SignIn(user.email, user.password)
		if err != nil {
			log.Fatalf("Error when login: %v\n", err)
		}

		req, err := http.NewRequest("PUT", "tasks/1/reminders/1", nil)
		if err != nil {
			t.Errorf("Error when getting reminders: %v\n", err)
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req = mux.SetURLVars(req, map[string]string{"task_id": "1", "reminder_id": "1"})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(s.UpdateReminder)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, http.StatusUnauthorized)
	}
}

func TestCreateReminder_WhenCreateReminderAsAdminOrCoach_401ErrorIsReturned(t *testing.T) {

	initializeTestData(t)

	users := []struct {
		email    string
		password     string
	}{
		{
			email: "raf@bot.com",
			password: "admin",
		},
		{
			email: "coach1@bot.com",
			password: "hellogophers",
		},
		{
			email: "coach2@bot.com",
			password: "hellogophers",
		},
	}

	for _, user := range users {
		token, err := s.SignIn(user.email, user.password)
		if err != nil {
			log.Fatalf("Error when login: %v\n", err)
		}

		req, err := http.NewRequest("POST", "tasks/1/reminders/", nil)
		if err != nil {
			t.Errorf("Error when getting reminders: %v\n", err)
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req = mux.SetURLVars(req, map[string]string{"task_id": "1"})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(s.CreateReminder)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, http.StatusUnauthorized)
	}
}

func TestDeleteReminder_WhenDeleteReminderAsAdminOrCoach_401ErrorIsReturned(t *testing.T) {

	initializeTestData(t)

	users := []struct {
		email    string
		password     string
	}{
		{
			email: "raf@bot.com",
			password: "admin",
		},
		{
			email: "coach1@bot.com",
			password: "hellogophers",
		},
		{
			email: "coach2@bot.com",
			password: "hellogophers",
		},
	}

	for _, user := range users {
		token, err := s.SignIn(user.email, user.password)
		if err != nil {
			log.Fatalf("Error when login: %v\n", err)
		}

		req, err := http.NewRequest("DELETE", "tasks/1/reminders/1", nil)
		if err != nil {
			t.Errorf("Error when getting reminders: %v\n", err)
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req = mux.SetURLVars(req, map[string]string{"task_id": "1", "reminder_id": "1"})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(s.GetReminders)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, http.StatusUnauthorized)
	}
}