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

func TestGetSessions_WhenGetSessionsAsAdmin_AllSessionsAreAvailable(t *testing.T) {

	initializeTestData(t)

	token, err := s.SignIn("raf@bot.com", "admin")
	if err != nil {
		log.Fatalf("Error when login: %v\n", err)
	}

	req, err := http.NewRequest("GET", "/sessions", nil)
	req.Header.Set("Authorization", "Bearer " + token)
	if err != nil {
		t.Errorf("Error when getting sessions: %v\n", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.GetSessions)
	handler.ServeHTTP(rr, req)

	var sessions []models.Session

	err = json.Unmarshal([]byte(rr.Body.String()), &sessions)
	if err != nil {
		log.Fatalf("Cannot convert to JSON: %v\n", err)
	}
	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(sessions), 4)
}

func TestGetSession_WhenGetSessionAsAdmin_SessionIsAvailable(t *testing.T) {

	initializeTestData(t)

	token, err := s.SignIn("raf@bot.com", "admin")
	if err != nil {
		log.Fatalf("Error when login: %v\n", err)
	}

	sessions := []struct {
		id 	string
		title string
		description string
		startingTime string
		dateSession string
	}{
		{
			id: "1",
			title: "First session Client 1",
			description: "Intro session: cardio/workout/abs",
			startingTime: "15:00",
			dateSession: "2022-11-21 15:00",
		},
		{
			id: "2",
			title: "Second session Client 1",
			description: "Follow-up session: cardio/workout/abs",
			startingTime: "15:00",
			dateSession: "2022-12-05 15:00",
		},
		{
			id: "3",
			title: "First session Client 2",
			description: "Intro session: cardio/workout/abs",
			startingTime: "15:00",
			dateSession: "2022-11-21 15:00",
		},
		{
			id: "4",
			title: "Second session Client 2",
			description: "Follow-up session: cardio/workout/abs",
			startingTime: "15:00",
			dateSession: "2022-12-05 15:00",
		},
	}

	for _, session := range sessions {

		req, err := http.NewRequest("GET", "/sessions/" + session.id, nil)
		req.Header.Set("Authorization", "Bearer " + token)
		if err != nil {
			t.Errorf("Error when getting session: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": session.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(s.GetSession)
		handler.ServeHTTP(rr, req)

		fmt.Println(rr.Body.String())

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to JSON: %v\n", err)
		}

		assert.Equal(t, rr.Code, 200)
		assert.Equal(t, session.title, responseMap["title"])
		assert.Equal(t, session.description, responseMap["description"])
		assert.Equal(t, session.startingTime, responseMap["starting_time"])
		assert.Equal(t, session.dateSession, responseMap["date_session (YYYY-MM-DD HH:MM)"])
	}
}

func TestGetSessions_WhenGetSessionsAsCoachOrClient_OnlySessionsInvolvingRequestUserAreAvailable(t *testing.T) {

	initializeTestData(t)

	users := []struct {
		email string
		userIdRequired string
		sessionsAvailable int
	}{
		{
			email: "coach1@bot.com",
			sessionsAvailable: 2,
		},
		{
			email: "client1@bot.com",
			sessionsAvailable: 2,
		},
		{
			email: "coach2@bot.com",
			sessionsAvailable: 2,
		},
		{
			email: "client2@bot.com",
			sessionsAvailable: 2,
		},
	}

	for _, user := range users {

		token, err := s.SignIn(user.email, "hellogophers")
		if err != nil {
			log.Fatalf("Error when login: %v\n", err)
		}

		req, err := http.NewRequest("GET", "/sessions", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		if err != nil {
			t.Errorf("Error when getting sessions: %v\n", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(s.GetSessions)
		handler.ServeHTTP(rr, req)

		fmt.Println(rr.Body.String())
		var sessions []models.Session

		err = json.Unmarshal([]byte(rr.Body.String()), &sessions)
		if err != nil {
			log.Fatalf("Cannot convert to JSON: %v\n", err)
		}
		assert.Equal(t, rr.Code, http.StatusOK)
		assert.Equal(t, len(sessions), 2)
	}
}

func TestGetSession_WhenGetSessionAsSessionCoachOrSessionClient_SessionIsAvailable(t *testing.T) {

	initializeTestData(t)

	users1 := []struct {
		email string
		userIdRequired string
	}{
		{
			email: "coach1@bot.com",
		},
		{
			email: "client1@bot.com",
		},
	}

	sessionsCoachAndClient1 := []struct {
		id 	string
		title string
		description string
		startingTime string
		dateSession string
	}{
		{
			id: "1",
			title: "First session Client 1",
			description: "Intro session: cardio/workout/abs",
			startingTime: "15:00",
			dateSession: "2022-11-21 15:00",
		},
		{
			id: "2",
			title: "Second session Client 1",
			description: "Follow-up session: cardio/workout/abs",
			startingTime: "15:00",
			dateSession: "2022-12-05 15:00",
		},

	}

	users2 := []struct {
		email string
		userIdRequired string
	}{
		{
			email: "coach2@bot.com",
		},
		{
			email: "client2@bot.com",
		},
	}

	sessionsCoachAndClient2 := []struct {
		id 	string
		title string
		description string
		startingTime string
		dateSession string
	}{

	{
	id: "3",
		title: "First session Client 2",
		description: "Intro session: cardio/workout/abs",
		startingTime: "15:00",
		dateSession: "2022-11-21 15:00",
	},
	{
	id: "4",
		title: "Second session Client 2",
		description: "Follow-up session: cardio/workout/abs",
		startingTime: "15:00",
		dateSession: "2022-12-05 15:00",
	},
	}

	for _, user := range users1 {

		token, err := s.SignIn(user.email, "hellogophers")
		if err != nil {
			log.Fatalf("Error when login: %v\n", err)
		}

		for _, session := range sessionsCoachAndClient1 {
				req, err := http.NewRequest("GET", "/sessions/" + session.id, nil)
				req.Header.Set("Authorization", "Bearer " + token)
				if err != nil {
					t.Errorf("Error when getting session: %v\n", err)
				}
				req = mux.SetURLVars(req, map[string]string{"id": session.id})
				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(s.GetSession)
				handler.ServeHTTP(rr, req)

				fmt.Println(rr.Body.String())

				responseMap := make(map[string]interface{})
				err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
				if err != nil {
					log.Fatalf("Cannot convert to JSON: %v\n", err)
				}

				assert.Equal(t, rr.Code, 200)
				assert.Equal(t, session.title, responseMap["title"])
				assert.Equal(t, session.description, responseMap["description"])
				assert.Equal(t, session.startingTime, responseMap["starting_time"])
				assert.Equal(t, session.dateSession, responseMap["date_session (YYYY-MM-DD HH:MM)"])
			}
		}


	for _, user := range users2 {

		token, err := s.SignIn(user.email, "hellogophers")
		if err != nil {
			log.Fatalf("Error when login: %v\n", err)
		}

		for _, session := range sessionsCoachAndClient2 {
			req, err := http.NewRequest("GET", "/sessions/" + session.id, nil)
			req.Header.Set("Authorization", "Bearer " + token)
			if err != nil {
				t.Errorf("Error when getting session: %v\n", err)
			}
			req = mux.SetURLVars(req, map[string]string{"id": session.id})
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(s.GetSession)
			handler.ServeHTTP(rr, req)

			fmt.Println(rr.Body.String())

			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
			if err != nil {
				log.Fatalf("Cannot convert to JSON: %v\n", err)
			}

			assert.Equal(t, rr.Code, 200)
			assert.Equal(t, session.title, responseMap["title"])
			assert.Equal(t, session.description, responseMap["description"])
			assert.Equal(t, session.startingTime, responseMap["starting_time"])
			assert.Equal(t, session.dateSession, responseMap["date_session (YYYY-MM-DD HH:MM)"])
		}
	}
}

func TestCreateSession_WhenCreateSessionAsAdmin_SessionIsCreated(t *testing.T) {

	initializeTestData(t)

	token, err := s.SignIn("raf@bot.com", "admin")
	if err != nil {
		log.Fatalf("Error when login: %v\n", err)
	}

	sessionsToCreate := []struct {
		inputJSON    string
		statusCode   int
		title string
		description string
		startingTime string
	}{
		{
			inputJSON:    `{"title": "Test session", "description":"Test session: cardio/workout/abs",
							"coach_id":2, "client_id":3, "year": 2022, "month":12, "day":9, 
							"starting_time":"16:01", "duration":60}`,
			statusCode:   201,
			title: "Test session",
			description: "Test session: cardio/workout/abs",
			startingTime: "16:01",
		},
		{
			inputJSON:    `{"title": "Test session", "description":"Test session: cardio/workout/abs",
							"coach_id":2, "client_id":3, "year": 2022, "month":12, "day":9, 
							"starting_time":"16:01", "duration":60}`,
			statusCode:   400,
			title: "Test session",
			description: "Test session: cardio/workout/abs",
			startingTime: "16:01",
		},
		{
			inputJSON:    `{"title": "Test session", "description":"Test session: cardio/workout/abs",
							"coach_id":2, "client_id":3, "year": 2022, "month":12, "day":9, 
							"starting_time":"18:00", "duration":60}`,
			statusCode:   201,
			title: "Test session",
			description: "Test session: cardio/workout/abs",
			startingTime: "18:00",
		},
	}

	for _, session := range sessionsToCreate {

		req, err := http.NewRequest("POST", "/sessions", bytes.NewBufferString(session.inputJSON))
		if err != nil {
			t.Errorf("Error when creating session: %v", err)
		}
		req.Header.Set("Authorization", "Bearer " + token)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(s.CreateSession)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, session.statusCode)

		if session.statusCode == 201 {
			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
			if err != nil {
				fmt.Printf("Cannot convert to JSON: %v", err)
			}

			assert.Equal(t, responseMap["title"], session.title)
			assert.Equal(t, responseMap["description"], session.description)
			assert.Equal(t, responseMap["starting_time"], session.startingTime)
		}
	}
}

func TestDeleteSession_WhenDeleteSessionAsAdmin_SessionIsDeleted(t *testing.T) {

	initializeTestData(t)

	token, err := s.SignIn("raf@bot.com", "admin")
	if err != nil {
		log.Fatalf("Error when login: %v\n", err)
	}

	req, err := http.NewRequest("DELETE", "/sessions", nil)
	if err != nil {
		t.Errorf("Error when deleting sessions: %v\n", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.DeleteSession)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusNoContent)

	req, err = http.NewRequest("GET", "/sessions", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	if err != nil {
		t.Errorf("Error when getting sessions: %v\n", err)
	}
	rr = httptest.NewRecorder()
	handler = s.GetSessions
	handler.ServeHTTP(rr, req)

	var sessions = []models.Session{}
	err = json.Unmarshal([]byte(rr.Body.String()), &sessions)
	if err != nil {
		log.Fatalf("Cannot convert to JSON: %v\n", err)
	}

	assert.Equal(t, len(sessions), 3)
}