package modeltests

import (
	"errors"
	"github.com/Alexandremerancienne/my_Sartorius/api/models"
	"gorm.io/gorm"
	"log"
	"testing"
)

func TestGetSessions(t *testing.T) {

	err := RefreshTables()
	if err != nil {
		t.Errorf("Error when refreshing tables: %v\n", err)
		return
	}

	_, _, _, _, _, _, err = SeedTables()
	if err != nil {
		log.Fatal(err)
	}

	u := models.Session{}
	allSessions, err := u.GetSessions(s.Database)

	if err != nil {
		t.Errorf("Error when getting Sessions: %v\n", err)
		return
	}
	if len(*allSessions) != 4 {
		t.Errorf("Number of Sessions should be equal to 4. Number of Sessions found: %d\n", len(*allSessions))
	}
}

func TestCreateSession(t *testing.T) {

	err := RefreshTables()
	if err != nil {
		t.Errorf("Error when refreshing tables: %v\n", err)
		return
	}

	_, _, _, _, _, _, err = SeedTables()
	if err != nil {
		log.Fatal(err)
	}

	var newSession models.Session
	newSession.Title = "Third session Client 1"
	newSession.Description =  "New session: burpees/push-ups/pull-ups"
	newSession.CoachID = 2
	newSession.ClientID = 3
	newSession.Year = 2022
	newSession.Month = 12
	newSession.Day = 07
	newSession.StartingTime = "15:00"
	newSession.Duration = 60
	newSession.DateSession = "2022-12-07 15:00"

	savedSession, err := newSession.CreateSession(s.Database)
	if err != nil {
		t.Errorf("Error when creating new session: %v\n", err)
		return
	}
	if newSession.ID != savedSession.ID {
		t.Errorf("savedSession ID should be equal to 5. ID registered: %d\n", savedSession.ID)
	}
	if newSession.Title != savedSession.Title {
		t.Errorf("savedSession title should be equal to 'First session Client 1'. Title registered: %s\n", savedSession.Title)
	}
	if newSession.Description != savedSession.Description {
		t.Errorf("savedSession description should be equal to 'Intro session: cardio/workout/abs'. Description registered: %s\n", savedSession.Description)
	}
	if newSession.CoachID != savedSession.CoachID {
		t.Errorf("savedSession coach ID should be equal to 2. Coach ID registered: %d\n", savedSession.CoachID)
	}
	if newSession.ClientID != savedSession.ClientID {
		t.Errorf("savedSession client ID should be equal to 3. Client ID registered: %d\n", savedSession.ClientID)
	}
	if newSession.Year != savedSession.Year {
		t.Errorf("savedSession year should be equal to 2022. Year registered: %d\n", savedSession.Year)
	}
	if newSession.Month != savedSession.Month {
		t.Errorf("savedSession month should be equal to 11. Month registered: %d\n", savedSession.Year)
	}
	if newSession.Day != savedSession.Day {
		t.Errorf("savedSession day should be equal to 07. Day registered: %d\n", savedSession.Day)
	}
	if newSession.StartingTime != savedSession.StartingTime {
		t.Errorf("savedSession starting time should be equal to 15:00. Starting time registered: %s\n", savedSession.StartingTime)
	}
	if newSession.Duration != savedSession.Duration {
		t.Errorf("savedSession duration should be equal to 60. Duration registered: %d\n", savedSession.Duration)
	}
	if newSession.DateSession != savedSession.DateSession {
		t.Errorf("savedSession date should be equal to '2022-12-07 15:00'. Session date registered: %s\n", savedSession.DateSession)
	}
}

func TestGetSession(t *testing.T) {

	err := RefreshTables()
	if err != nil {
		t.Errorf("Error when refreshing tables: %v\n", err)
		return
	}

	_, _, _, _, _, _, err = SeedTables()
	if err != nil {
		log.Fatal(err)
	}

	u := models.Session{}

	retrievedSession, err := u.GetSessionByID(s.Database, 1)
	if err != nil {
		t.Errorf("Error when getting session: %v\n", err)
		return
	}
	if retrievedSession.ID != 1 {
		t.Errorf("savedSession ID should be equal to 1. ID registered: %d\n", retrievedSession.ID)
	}
	if retrievedSession.Title != "First session Client 1" {
		t.Errorf("savedSession title should be equal to 'First session Client 1'. Title registered: %s\n", retrievedSession.Title)
	}
	if retrievedSession.Description != "Intro session: cardio/workout/abs" {
		t.Errorf("savedSession description should be equal to 'Intro session: cardio/workout/abs'. Description registered: %s\n", retrievedSession.Description)
	}
	if retrievedSession.CoachID != 2 {
		t.Errorf("savedSession coach ID should be equal to 2. Coach ID registered: %d\n", retrievedSession.CoachID)
	}
	if retrievedSession.ClientID != 3 {
		t.Errorf("savedSession client ID should be equal to 3. Client ID registered: %d\n", retrievedSession.ClientID)
	}
	if retrievedSession.Year != 2022 {
		t.Errorf("savedSession year should be equal to 2022. Year registered: %d\n", retrievedSession.Year)
	}
	if retrievedSession.Month != 11 {
		t.Errorf("savedSession month should be equal to 11. Month registered: %d\n", retrievedSession.Year)
	}
	if retrievedSession.Day != 21 {
		t.Errorf("savedSession day should be equal to 21. Day registered: %d\n", retrievedSession.Day)
	}
	if retrievedSession.StartingTime != "15:00" {
		t.Errorf("savedSession starting time should be equal to 15:00. Starting time registered: %s\n", retrievedSession.StartingTime)
	}
	if retrievedSession.Duration != 60 {
		t.Errorf("savedSession duration should be equal to 60. Duration registered: %d\n", retrievedSession.Duration)
	}
	if retrievedSession.DateSession != "2022-11-21 15:00" {
		t.Errorf("savedSession date should be equal to '2022-11-21 15:00'. Session date registered: %s\n", retrievedSession.DateSession)
	}
}

func TestDeleteSession(t *testing.T) {

	err := RefreshTables()
	if err != nil {
		t.Errorf("Error when refreshing tables: %v\n", err)
		return
	}

	_, _, _, _, _, _, err = SeedTables()
	if err != nil {
		log.Fatal(err)
	}

	u := models.Session{}

	_, err = u.DeleteSession(s.Database, 4)
	if err != nil {
		t.Errorf("Error when deleting session: %v\n", err)
		return
	}

	i := models.Session{}
	allSessions, err := i.GetSessions(s.Database)

	if err != nil {
		t.Errorf("Error when getting sessions: %v\n", err)
		return
	}

	if len(*allSessions) != 3 {
		t.Errorf("Number of sessions should be equal to 4. Number of sessions found: %d\n", len(*allSessions))
	}

	j := models.Session{}

	if retrievedSession, err := j.GetSessionByID(s.Database, 4) ; !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("Session with ID 4 should have been deleted. One Session found with ID 4: %v\n", retrievedSession)
	}
}