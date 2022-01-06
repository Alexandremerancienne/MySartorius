package modeltests

import (
	"errors"
	"github.com/Alexandremerancienne/my_Sartorius/api/models"
	"gorm.io/gorm"
	"log"
	"testing"
)

func TestGetTasks(t *testing.T) {

	err := RefreshTables()
	if err != nil {
		t.Errorf("Error when refreshing tables: %v\n", err)
		return
	}

	_, _, _, _, _, _, err = SeedTables()
	if err != nil {
		log.Fatal(err)
	}

	u := models.Task{}
	allTasks, err := u.GetTasks(s.Database)

	if err != nil {
		t.Errorf("Error when getting tasks: %v\n", err)
		return
	}
	if len(*allTasks) != 4 {
		t.Errorf("Number of Tasks should be equal to 4. Number of Tasks found: %d\n", len(*allTasks))
	}
}

func TestCreateTask(t *testing.T) {

	err := RefreshTables()
	if err != nil {
		t.Errorf("Error when refreshing tables: %v\n", err)
		return
	}

	_, _, _, _, _, _, err = SeedTables()
	if err != nil {
		log.Fatal(err)
	}


	var newTask models.Task

	newTask.Title = "Friday workout"
	newTask.Description = "10 abs, 10 push-ups, 10 burpees. 40 seconds rest between each exercise"
	newTask.AssignerID = 4
	newTask.AssigneeID = 5
	newTask.Year = 2022
	newTask.Month = 11
	newTask.Day = 19
	newTask.Duration = 30
	newTask.DateTask = "2022-11-19 15:00"


	savedTask, err := newTask.CreateTask(s.Database)
	if err != nil {
		t.Errorf("Error when creating new Task: %v\n", err)
		return
	}
	if newTask.ID != savedTask.ID {
		t.Errorf("savedTask ID should be equal to 5. ID registered: %d\n", savedTask.ID)
	}
	if newTask.Title != savedTask.Title {
		t.Errorf("savedTask title should be equal to 'Friday workout'. Title registered: %s\n", savedTask.Title)
	}
	if newTask.Description != savedTask.Description {
		t.Errorf("savedTask description should be equal to '10 abs, 10 push-ups, 10 burpees. 40 seconds rest between each exercise'. Description registered: %s\n", savedTask.Description)
	}
	if newTask.AssignerID != savedTask.AssignerID {
		t.Errorf("savedTask assigner ID should be equal to 4. Assigner ID registered: %d\n", savedTask.AssignerID)
	}
	if newTask.AssigneeID != savedTask.AssigneeID {
		t.Errorf("savedTask assignee ID should be equal to 5. Assignee ID registered: %d\n", savedTask.AssigneeID)
	}
	if newTask.Year != savedTask.Year {
		t.Errorf("savedTask year should be equal to 2022. Year registered: %d\n", savedTask.Year)
	}
	if newTask.Month != savedTask.Month {
		t.Errorf("savedTask month should be equal to 11. Month registered: %d\n", savedTask.Year)
	}
	if newTask.Day != savedTask.Day {
		t.Errorf("savedTask day should be equal to 19. Day registered: %d\n", savedTask.Day)
	}
	if newTask.Duration != savedTask.Duration {
		t.Errorf("savedTask duration should be equal to 30. Duration registered: %d\n", savedTask.Duration)
	}
	if newTask.DateTask != savedTask.DateTask {
		t.Errorf("savedTask date should be equal to '2022-11-19 15:00'. Task date registered: %s\n", savedTask.DateTask)
	}
}


func TestGetTask(t *testing.T) {

	err := RefreshTables()
	if err != nil {
		t.Errorf("Error when refreshing tables: %v\n", err)
		return
	}

	_, _, _, _, _, _, err = SeedTables()
	if err != nil {
		log.Fatal(err)
	}

	u := models.Task{}

	retrievedTask, err := u.GetTaskByID(s.Database, 4)
	if err != nil {
		t.Errorf("Error when getting Task: %v\n", err)
		return
	}
	if retrievedTask.ID != 4 {
		t.Errorf("savedTask ID should be equal to 4. ID registered: %d\n", retrievedTask.ID)
	}
	if retrievedTask.Title != "Thursday workout" {
		t.Errorf("savedTask title should be equal to 'Thursday workout'. Title registered: %s\n", retrievedTask.Title)
	}
	if retrievedTask.Description != "7 abs, 8 push-ups, 8 burpees. 45 seconds rest between each exercise" {
		t.Errorf("savedTask description should be equal to `7 abs, 8 push-ups, 8 burpees. 45 seconds rest between each exercise`. Description registered: %s\n", retrievedTask.Description)
	}
	if retrievedTask.AssignerID != 4 {
		t.Errorf("savedTask assigner ID should be equal to 4. Assigner ID registered: %d\n", retrievedTask.AssignerID)
	}
	if retrievedTask.AssigneeID != 5 {
		t.Errorf("savedTask assignee ID should be equal to 5. Assignee ID registered: %d\n", retrievedTask.AssigneeID)
	}
	if retrievedTask.Year != 2022 {
		t.Errorf("savedTask year should be equal to 2022. Year registered: %d\n", retrievedTask.Year)
	}
	if retrievedTask.Month != 11 {
		t.Errorf("savedTask month should be equal to 11. Month registered: %d\n", retrievedTask.Year)
	}
	if retrievedTask.Day != 18 {
		t.Errorf("savedTask day should be equal to 18. Day registered: %d\n", retrievedTask.Day)
	}
	if retrievedTask.Duration != 30 {
		t.Errorf("savedTask duration should be equal to 30. Duration registered: %d\n", retrievedTask.Duration)
	}
	if retrievedTask.DateTask != "2022-11-18" {
		t.Errorf("savedTask date should be equal to '2022-11-18'. Task date registered: %s\n", retrievedTask.DateTask)
	}
}

func TestDeleteTask(t *testing.T) {

	err := RefreshTables()
	if err != nil {
		t.Errorf("Error when refreshing tables: %v\n", err)
		return
	}

	_, _, _, _, _, _, err = SeedTables()
	if err != nil {
		log.Fatal(err)
	}

	u := models.Task{}

	_, err = u.DeleteTask(s.Database, 4)
	if err != nil {
		t.Errorf("Error when deleting Task: %v\n", err)
		return
	}

	i := models.Task{}
	allTasks, err := i.GetTasks(s.Database)

	if err != nil {
		t.Errorf("Error when getting Tasks: %v\n", err)
		return
	}

	if len(*allTasks) != 3 {
		t.Errorf("Number of Tasks should be equal to 4. Number of Tasks found: %d\n", len(*allTasks))
	}

	j := models.Task{}

	if retrievedTask, err := j.GetTaskByID(s.Database, 4) ; !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("Task with ID 4 should have been deleted. One Task found with ID 4: %v\n", retrievedTask)
	}
}