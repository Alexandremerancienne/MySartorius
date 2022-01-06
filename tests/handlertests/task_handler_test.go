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

func TestGetTasks_WhenGetTasksAsAdmin_AllTasksAreAvailable(t *testing.T) {

	initializeTestData(t)

	token, err := s.SignIn("raf@bot.com", "admin")
	if err != nil {
		log.Fatalf("Error when login: %v\n", err)
	}

	req, err := http.NewRequest("GET", "/tasks", nil)
	req.Header.Set("Authorization", "Bearer " + token)
	if err != nil {
		t.Errorf("Error when getting tasks: %v\n", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.GetTasks)
	handler.ServeHTTP(rr, req)

	var tasks []models.Task

	err = json.Unmarshal([]byte(rr.Body.String()), &tasks)
	if err != nil {
		log.Fatalf("Cannot convert to JSON: %v\n", err)
	}
	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(tasks), 4)
}

func TestGetTask_WhenGetTaskAsAdmin_TaskIsAvailable(t *testing.T) {

	initializeTestData(t)

	token, err := s.SignIn("raf@bot.com", "admin")
	if err != nil {
		log.Fatalf("Error when login: %v\n", err)
	}

	tasks := []struct {
		id 	string
		title string
		description string
		dateTask string
	}{
		{
			id: "1",
			title: "Monday workout",
			description: "10 abs, 15 push-ups, 10 burpees. 30 seconds rest between each exercise",
			dateTask: "2022-11-15",
		},
		{
			id: "2",
			title: "Wednesday workout",
			description: "15 abs, 20 push-ups, 15 burpees. 35 seconds rest between each exercise",
			dateTask: "2022-11-17",
		},
		{
			id: "3",
			title: "Tuesday workout",
			description: "5 abs, 5 push-ups, 5 burpees. 40 seconds rest between each exercise",
			dateTask: "2022-11-16",
		},
		{
			id: "4",
			title: "Thursday workout",
			description: "7 abs, 8 push-ups, 8 burpees. 45 seconds rest between each exercise",
			dateTask: "2022-11-18",
		},
	}

	for _, task := range tasks {

		req, err := http.NewRequest("GET", "/tasks/" + task.id, nil)
		req.Header.Set("Authorization", "Bearer " + token)
		if err != nil {
			t.Errorf("Error when getting task: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": task.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(s.GetTask)
		handler.ServeHTTP(rr, req)

		fmt.Println(rr.Body.String())

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to JSON: %v\n", err)
		}

		assert.Equal(t, rr.Code, 200)
		assert.Equal(t, task.title, responseMap["title"])
		assert.Equal(t, task.description, responseMap["description"])
		assert.Equal(t, task.dateTask, responseMap["date_task (YYYY-MM-DD)"])
	}
}

func TestGetTasks_WhenGetTasksAsCoachOrClient_OnlyTasksInvolvingRequestUserAreAvailable(t *testing.T) {

	initializeTestData(t)

	users := []struct {
		email string
		userIdRequired string
		tasksAvailable int
	}{
		{
			email: "coach1@bot.com",
			tasksAvailable: 2,
		},
		{
			email: "client1@bot.com",
			tasksAvailable: 2,
		},
		{
			email: "coach2@bot.com",
			tasksAvailable: 2,
		},
		{
			email: "client2@bot.com",
			tasksAvailable: 2,
		},
	}

	for _, user := range users {

		token, err := s.SignIn(user.email, "hellogophers")
		if err != nil {
			log.Fatalf("Error when login: %v\n", err)
		}

		req, err := http.NewRequest("GET", "/tasks", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		if err != nil {
			t.Errorf("Error when getting tasks: %v\n", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(s.GetTasks)
		handler.ServeHTTP(rr, req)

		fmt.Println(rr.Body.String())
		var tasks []models.Task

		err = json.Unmarshal([]byte(rr.Body.String()), &tasks)
		if err != nil {
			log.Fatalf("Cannot convert to JSON: %v\n", err)
		}
		assert.Equal(t, rr.Code, http.StatusOK)
		assert.Equal(t, len(tasks), 2)
	}
}

func TestGetTask_WhenGetTaskAsTaskAssignerOrTaskAssignee_TaskIsAvailable(t *testing.T) {

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

	tasksCoachAndClient1 := []struct {
		id 	string
		title string
		description string
		dateTask string
	}{
		{
			id: "1",
			title: "Monday workout",
			description: "10 abs, 15 push-ups, 10 burpees. 30 seconds rest between each exercise",
			dateTask: "2022-11-15",
		},
		{
			id: "2",
			title: "Wednesday workout",
			description: "15 abs, 20 push-ups, 15 burpees. 35 seconds rest between each exercise",
			dateTask: "2022-11-17",
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

	tasksCoachAndClient2 := []struct {
		id 	string
		title string
		description string
		dateTask string
	}{
		{
			id: "3",
			title: "Tuesday workout",
			description: "5 abs, 5 push-ups, 5 burpees. 40 seconds rest between each exercise",
			dateTask: "2022-11-16",
		},
		{
			id: "4",
			title: "Thursday workout",
			description: "7 abs, 8 push-ups, 8 burpees. 45 seconds rest between each exercise",
			dateTask: "2022-11-18",
		},
	}

	for _, user := range users1 {

		token, err := s.SignIn(user.email, "hellogophers")
		if err != nil {
			log.Fatalf("Error when login: %v\n", err)
		}

		for _, task := range tasksCoachAndClient1 {
			req, err := http.NewRequest("GET", "/tasks/" + task.id, nil)
			req.Header.Set("Authorization", "Bearer " + token)
			if err != nil {
				t.Errorf("Error when getting task: %v\n", err)
			}
			req = mux.SetURLVars(req, map[string]string{"id": task.id})
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(s.GetTask)
			handler.ServeHTTP(rr, req)

			fmt.Println(rr.Body.String())

			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
			if err != nil {
				log.Fatalf("Cannot convert to JSON: %v\n", err)
			}

			assert.Equal(t, rr.Code, 200)
			assert.Equal(t, task.title, responseMap["title"])
			assert.Equal(t, task.description, responseMap["description"])
			assert.Equal(t, task.dateTask, responseMap["date_task (YYYY-MM-DD)"])
		}
	}

	for _, user := range users2 {

		token, err := s.SignIn(user.email, "hellogophers")
		if err != nil {
			log.Fatalf("Error when login: %v\n", err)
		}

		for _, task := range tasksCoachAndClient2 {
			req, err := http.NewRequest("GET", "/tasks/" + task.id, nil)
			req.Header.Set("Authorization", "Bearer " + token)
			if err != nil {
				t.Errorf("Error when getting task: %v\n", err)
			}
			req = mux.SetURLVars(req, map[string]string{"id": task.id})
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(s.GetTask)
			handler.ServeHTTP(rr, req)

			fmt.Println(rr.Body.String())

			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
			if err != nil {
				log.Fatalf("Cannot convert to JSON: %v\n", err)
			}

			assert.Equal(t, rr.Code, 200)
			assert.Equal(t, task.title, responseMap["title"])
			assert.Equal(t, task.description, responseMap["description"])
			assert.Equal(t, task.dateTask, responseMap["date_task (YYYY-MM-DD)"])
		}
	}
}

func TestCreateTask_WhenCreateTaskAsAdmin_TaskIsCreated(t *testing.T) {

	initializeTestData(t)

	token, err := s.SignIn("raf@bot.com", "admin")
	if err != nil {
		log.Fatalf("Error when login: %v\n", err)
	}

	tasksToCreate := []struct {
		inputJSON    string
		statusCode   int
		title string
		description string
	}{
		{
			inputJSON:    `{"title": "Monday workout", 
							"description":"10 abs, 15 push-ups, 10 burpees. 30 seconds rest between each exercise",
							"assigner_id": 2, "assignee_id": 3, "year": 2021, "month":11, "day":20, "duration":60}`,
			statusCode:   400,
			title: "Monday workout",
			description: "10 abs, 15 push-ups, 10 burpees. 30 seconds rest between each exercise",
		},
		{
			inputJSON:    `{"title": "Monday workout", 
							"description":"10 abs, 15 push-ups, 10 burpees. 30 seconds rest between each exercise",
							"assigner_id": 2, "assignee_id": 3, "year": 2022, "month":11, "day":20, "duration":60}`,
			statusCode:   201,
			title: "Monday workout",
			description: "10 abs, 15 push-ups, 10 burpees. 30 seconds rest between each exercise",
		},
		{
			inputJSON:    `{"title": "Monday workout (2nd task)", 
							"description":"15 abs, 20 push-ups, 11 burpees. 30 seconds rest between each exercise",
							"assigner_id": 2, "assignee_id": 3, "year": 2022, "month":11, "day":20, "duration":60}`,
			statusCode:   201,
			title: "Monday workout (2nd task)",
			description: "15 abs, 20 push-ups, 11 burpees. 30 seconds rest between each exercise",
		},
		{
			inputJSON:    `{"title": "Tuesday workout", 
							"description":"11 abs, 11 push-ups, 5 burpees. 30 seconds rest between each exercise",
							"assigner_id": 2, "assignee_id": 3, "year": 2022, "month":11, "day":21, "duration":60}`,
			statusCode:   201,
			title: "Tuesday workout",
			description: "11 abs, 11 push-ups, 5 burpees. 30 seconds rest between each exercise",
		},
	}

	for _, task := range tasksToCreate {

		req, err := http.NewRequest("POST", "/tasks", bytes.NewBufferString(task.inputJSON))
		if err != nil {
			t.Errorf("Error when creating task: %v", err)
		}
		req.Header.Set("Authorization", "Bearer " + token)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(s.CreateTask)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, task.statusCode)

		if task.statusCode == 201 {
			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
			if err != nil {
				fmt.Printf("Cannot convert to JSON: %v", err)
			}

			assert.Equal(t, responseMap["title"], task.title)
			assert.Equal(t, responseMap["description"], task.description)
		}
	}
}

func TestDeleteTask_WhenDeleteTaskAsAdmin_TaskIsDeleted(t *testing.T) {

	initializeTestData(t)

	token, err := s.SignIn("raf@bot.com", "admin")
	if err != nil {
		log.Fatalf("Error when login: %v\n", err)
	}

	req, err := http.NewRequest("DELETE", "/tasks", nil)
	if err != nil {
		t.Errorf("Error when deleting tasks: %v\n", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.DeleteTask)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusNoContent)

	req, err = http.NewRequest("GET", "/tasks", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	if err != nil {
		t.Errorf("Error when getting tasks: %v\n", err)
	}
	rr = httptest.NewRecorder()
	handler = s.GetTasks
	handler.ServeHTTP(rr, req)

	var tasks = []models.Task{}
	err = json.Unmarshal([]byte(rr.Body.String()), &tasks)
	if err != nil {
		log.Fatalf("Cannot convert to JSON: %v\n", err)
	}

	assert.Equal(t, len(tasks), 3)
}