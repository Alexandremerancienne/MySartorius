package modeltests

import (
	"errors"
	"github.com/Alexandremerancienne/my_Sartorius/api/models"
	"gorm.io/gorm"
	"log"
	"testing"
)

var err error

func TestGetUsers(t *testing.T) {

	err := RefreshTables()
	if err != nil {
		t.Errorf("Error when refreshing tables: %v\n", err)
		return
	}

	_, _, _, _, _, _, err = SeedTables()
	if err != nil {
		log.Fatal(err)
	}

	u := models.User{}
	allUsers, err := u.GetUsers(s.Database)

	if err != nil {
		t.Errorf("Error when getting users: %v\n", err)
		return
	}
	if len(*allUsers) != 5 {
		t.Errorf("Number of users should be equal to 5. Number of users found: %d\n", len(*allUsers))
	}
}

func TestCreateUser(t *testing.T) {

	err := RefreshTables()
	if err != nil {
		t.Errorf("Error when refreshing tables: %v\n", err)
		return
	}

	var newUser models.User
	newUser.FirstName = "Coach 1"
	newUser.LastName =    "Whatever"
	newUser.Email = "coach1@bot.com"
	newUser.Password = "hellogophers"
	newUser.PhoneNumber =    "+33646452149"
	newUser.Role = "coach"

	savedUser, err := newUser.CreateUser(s.Database)
	if err != nil {
		t.Errorf("Error when creating new user: %v\n", err)
		return
	}
	if newUser.ID != savedUser.ID {
		t.Errorf("savedUser ID should be equal to 1. ID registered: %d\n", savedUser.ID)
	}
	if newUser.FirstName != savedUser.FirstName {
		t.Errorf("savedUser first name should be equal to Coach1. First Name registered: %s\n", savedUser.FirstName)
	}
	if newUser.Email != savedUser.Email {
		t.Errorf("savedUser email should be equal to coach1@bot.com. Email registered: %s\n", savedUser.Email)
	}
	if newUser.Password != savedUser.Password {
		t.Errorf("Incorrect password when creating new user\n")
	}
	if newUser.PhoneNumber != savedUser.PhoneNumber {
		t.Errorf("savedUser phone number should be equal to +33646452149. Number registered: %s\n", savedUser.PhoneNumber)
	}
	if newUser.Role != savedUser.Role {
		t.Errorf("savedUser role should be equal to coach. Role registered: %s\n", savedUser.Role)
	}
}

func TestGetUser(t *testing.T) {

	err := RefreshTables()
	if err != nil {
		t.Errorf("Error when refreshing tables: %v\n", err)
		return
	}

	_, _, _, _, _, _, err = SeedTables()
	if err != nil {
		log.Fatal(err)
	}

	u := models.User{}

	retrievedUser, err := u.GetUserByID(s.Database, 1)
	if err != nil {
		t.Errorf("Error when getting user: %v\n", err)
		return
	}

	if retrievedUser.ID != 1 {
		t.Errorf("retrievedUser ID should be equal to 1. ID registered: %d\n", retrievedUser.ID)
	}
	if retrievedUser.FirstName != "Raphaël" {
		t.Errorf("retrievedUser first name should be equal to Raphaël. First Name registered: %s\n", retrievedUser.FirstName)
	}
	if retrievedUser.LastName != "Merancienne" {
		t.Errorf("retrievedUser last name should be equal to Merancienne. Last Name registered: %s\n", retrievedUser.LastName)
	}
	if retrievedUser.Email != "raf@bot.com" {
		t.Errorf("retrievedUser last name should be equal to coach1@bot.com. Email registered: %s\n", retrievedUser.Email)
	}
	if retrievedUser.PhoneNumber != "+33497854010" {
		t.Errorf("retrievedUser phone number should be equal to +33646452149. Number registered: %s\n", retrievedUser.PhoneNumber)
	}
	if retrievedUser.Role != "manager" {
		t.Errorf("retrievedUser role should be equal to coach. Role registered: %s\n", retrievedUser.Role)
	}
}

func TestDeleteUser(t *testing.T) {

	err := RefreshTables()
	if err != nil {
		t.Errorf("Error when refreshing tables: %v\n", err)
		return
	}

	_, _, _, _, _, _, err = SeedTables()
	if err != nil {
		log.Fatal(err)
	}

	u := models.User{}

	_, err = u.DeleteUser(s.Database, 5)
	if err != nil {
		t.Errorf("Error when deleting user: %v\n", err)
		return
	}

	i := models.User{}
	allUsers, err := i.GetUsers(s.Database)

	if err != nil {
		t.Errorf("Error when getting users: %v\n", err)
		return
	}

	if len(*allUsers) != 4 {
		t.Errorf("Number of users should be equal to 4. Number of users found: %d\n", len(*allUsers))
	}

	j := models.User{}

	if retrievedUser, err := j.GetUserByID(s.Database, 5) ; !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("User with ID 5 should have been deleted. One user found with ID 5: %v\n", retrievedUser)
	}
}