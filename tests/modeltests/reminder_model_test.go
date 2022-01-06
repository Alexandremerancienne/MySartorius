package modeltests

import (
	"errors"
	"github.com/Alexandremerancienne/my_Sartorius/api/models"
	"gorm.io/gorm"
	"log"
	"testing"
)

func TestGetReminders(t *testing.T) {

	err := RefreshTables()
	if err != nil {
		t.Errorf("Error when refreshing tables: %v\n", err)
		return
	}

	_, _, _, _, _, _, err = SeedTables()
	if err != nil {
		log.Fatal(err)
	}

	u := models.Reminder{}
	allReminders, err := u.GetReminders(s.Database, 4)

	if err != nil {
		t.Errorf("Error when getting reminders: %v\n", err)
		return
	}
	if len(*allReminders) != 2 {
		t.Errorf("Number of Reminders should be equal to 2. Number of Reminders found: %d\n", len(*allReminders))
	}
}

func TestCreateReminder(t *testing.T) {

	err := RefreshTables()
	if err != nil {
		t.Errorf("Error when refreshing tables: %v\n", err)
		return
	}

	_, _, _, _, _, _, err = SeedTables()
	if err != nil {
		log.Fatal(err)
	}

	var newReminder models.Reminder

	newReminder.TaskID= 1
	newReminder.Description = "Take a deep inspiration between abs"

	savedReminder, err := newReminder.CreateReminder(s.Database, 1)
	if err != nil {
		t.Errorf("Error when creating new reminder: %v\n", err)
		return
	}
	if newReminder.ID != savedReminder.ID {
		t.Errorf("savedReminder ID should be equal to 4. ID registered: %d\n", savedReminder.ID)
	}
	if newReminder.TaskID != savedReminder.TaskID {
		t.Errorf("savedReminder task ID should be equal to 1. Task ID registered: %d\n", savedReminder.TaskID)
	}
	if newReminder.Description != savedReminder.Description {
		t.Errorf("savedReminder description should be equal to 'Take a deep inspiration between abs'. Description registered: %s\n", savedReminder.Description)
	}
}


func TestGetReminder(t *testing.T) {

	err := RefreshTables()
	if err != nil {
		t.Errorf("Error when refreshing tables: %v\n", err)
		return
	}

	_, _, _, _, _, _, err = SeedTables()
	if err != nil {
		log.Fatal(err)
	}

	u := models.Reminder{}

	retrievedReminder, err := u.GetReminderByID(s.Database, 1, 1)
	if err != nil {
		t.Errorf("Error when getting reminder: %v\n", err)
		return
	}

	if retrievedReminder.ID != 1 {
		t.Errorf("savedReminder ID should be equal to 4. ID registered: %d\n", retrievedReminder.ID)
	}
	if retrievedReminder.TaskID != 1 {
		t.Errorf("savedReminder task ID should be equal to 1. Task ID registered: %d\n", retrievedReminder.TaskID)
	}
	if retrievedReminder.Description != "Do not forget to bend the arms on the push-ups" {
		t.Errorf("savedReminder description should be equal to 'Do not forget to bend the arms on the push-ups'. Description registered: %s\n", retrievedReminder.Description)
	}
}

func TestDeleteReminder(t *testing.T) {

	err := RefreshTables()
	if err != nil {
		t.Errorf("Error when refreshing tables: %v\n", err)
		return
	}

	_, _, _, _, _, _, err = SeedTables()
	if err != nil {
		log.Fatal(err)
	}

	u := models.Reminder{}

	_, err = u.DeleteReminder(s.Database, 1, 4)
	if err != nil {
		t.Errorf("Error when deleting reminder: %v\n", err)
		return
	}

	i := models.Reminder{}
	allReminders, err := i.GetReminders(s.Database, 1)

	if err != nil {
		t.Errorf("Error when getting reminders: %v\n", err)
		return
	}

	if len(*allReminders) != 1 {
		t.Errorf("Number of reminders for task 1 should be equal to 1. Number of Reminders found: %d\n", len(*allReminders))
	}

	j := models.Reminder{}

	if retrievedReminder, err := j.GetReminderByID(s.Database, 1, 4) ; !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("Reminder with ID 4 for task 1 should have been deleted. One Reminder found with ID 4 for task 1: %v\n", retrievedReminder)
	}
}