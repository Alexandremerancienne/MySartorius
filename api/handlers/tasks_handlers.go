package handlers

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/Alexandremerancienne/my_Sartorius/api/auth"
	"github.com/Alexandremerancienne/my_Sartorius/api/exceptions"

	"errors"
	//"github.com/Alexandremerancienne/my_Sartorius/api/auth"
	//"github.com/Alexandremerancienne/my_Sartorius/api/exceptions"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Alexandremerancienne/my_Sartorius/api/models"
	"github.com/gorilla/mux"
)

func (server *Server) CreateTask(w http.ResponseWriter, r *http.Request) {

	user := models.User{}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		err = errors.New("Invalid token: please login with valid credentials.")
		exceptions.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	k, err := user.GetUserByID(server.Database, tokenID)
	if err != nil {
		err = errors.New("500 Internal Server Error: cannot get user")
		exceptions.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if tokenID == uint32(k.ID) &&
		(k.Role == "manager" || k.Role == "coach") {

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			exceptions.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}
		task := models.Task{}
		if k.Role == "coach" {
			task.AssignerID = int(k.ID)
		}
		if err = json.Unmarshal(body, &task); err != nil {
			err = errors.New("422 Unprocessable Entity Error: Cannot parse JSON data.")
			exceptions.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}

		taskMonth := ""
		if len(strconv.Itoa(task.Month)) == 1 {
			taskMonth = "0" + strconv.Itoa(task.Month)
		} else {
			taskMonth = strconv.Itoa(task.Month)
		}

		taskDay := ""
		if len(strconv.Itoa(task.Day)) == 1 {
			taskDay = "0" + strconv.Itoa(task.Day)
		} else {
			taskDay = strconv.Itoa(task.Day)
		}

		t1, err := time.Parse(
			time.RFC3339,
			strconv.Itoa(task.Year)+"-"+taskMonth+"-"+taskDay+"T"+"00:00"+":00+00:00")

		if err != nil {
			err = errors.New("400 Bad Request Error: Cannot create task with this date")
			exceptions.ERROR(w, http.StatusBadRequest, err)
			return
		}

		if t1.Format(time.RFC3339) < time.Now().Format(time.RFC3339) {
			err = errors.New("400 Bad Request Error: Cannot create task with a past date")
			exceptions.ERROR(w, http.StatusBadRequest, err)
			return
		}

		task.DateTask = strings.Replace(t1.Format(time.RFC3339)[:10], "T", " ", -1)

		newTask, err := task.CreateTask(server.Database)

		if err != nil {
			err = errors.New("Error: A task cannot be created for these users. Please check the IDs Provided.")
			exceptions.ERROR(w, http.StatusBadRequest, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(newTask)

	} else {
		err = errors.New("Missing credentials: access restricted to Management and coaches.")
		exceptions.ERROR(w, http.StatusUnauthorized, err)
		return
	}
}

func (server *Server) GetTasks(w http.ResponseWriter, r *http.Request) {

	u, err := url.Parse(r.URL.String())
	if err != nil {
		err = errors.New("422 Unprocessable Entity Error: Cannot parse url")
		exceptions.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	qs := u.Query()
	y := server.readInt(qs, "coach_id", 0)
	z := server.readInt(qs, "client_id", 0)

	tasks := models.Task{}
	retrievedTasks, err := tasks.GetTasks(server.Database)
	if err != nil {
		err = errors.New("500 Internal Server Error: cannot get tasks")
		exceptions.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	tasksSlice := []models.Task{}

	for _, t := range *retrievedTasks {

		u := models.User{}
		retrievedCoach, err := u.GetCoachByID(server.Database, t.AssignerID)

		j := models.User{}
		retrievedClient, err := j.GetClientByID(server.Database, t.AssigneeID)

		user := models.User{}
		tokenID, err := auth.ExtractTokenID(r)
		if err != nil {
			err = errors.New("Invalid token: please login with valid credentials.")
			exceptions.ERROR(w, http.StatusUnauthorized, err)
			return
		}

		k, err := user.GetUserByID(server.Database, tokenID)
		if err != nil {
			err = errors.New("500 Internal Server Error: cannot get user")
			exceptions.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		if tokenID == uint32(k.ID) &&
			(k.Role == "manager" ||
				(k.Role == "coach" && k.ID == retrievedCoach.ID) ||
				(k.Role == "client" && k.ID == retrievedClient.ID)) {

			tasksSlice = append(tasksSlice, t)
		} else {
			continue
		}
	}

	t := models.Task{}

	if y != 0 || z != 0 {
		w.Header().Set("Content-Type", "application/json")
		tasks, err := t.GetTasks(server.Database)
		if err != nil {
			err = errors.New("500 Internal Server Error: cannot get tasks")
			exceptions.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		q := models.User{}
		tokenID, err := auth.ExtractTokenID(r)
		if err != nil {
			err = errors.New("Invalid token: Please login with valid credentials.")
			exceptions.ERROR(w, http.StatusUnauthorized, err)
			return
		}

		j, err := q.GetUserByID(server.Database, tokenID)
		if err != nil {
			err = errors.New("500 Internal Server Error: Cannot get user")
			exceptions.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		server.ApplyTaskFilters(w, tasks, y, z, *j)
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasksSlice)
	}
}


func (server *Server) GetTask(w http.ResponseWriter, r *http.Request) {

	task := models.Task{}
	params := mux.Vars(r)
	taskID, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		err = errors.New("422 Unprocessable Entity Error: Cannot parse url")
		exceptions.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	retrievedTask, err := task.GetTaskByID(server.Database, uint32(taskID))

	u := models.User{}
	retrievedCoach, err := u.GetCoachByID(server.Database, retrievedTask.AssignerID)

	u2 := models.User{}
	retrievedClient, err := u2.GetClientByID(server.Database, retrievedTask.AssigneeID)

	user := models.User{}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		err = errors.New("Invalid token: please login with valid credentials.")
		exceptions.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	k, err := user.GetUserByID(server.Database, tokenID)
	if err != nil {
		err = errors.New("500 Internal Server Error: cannot get user")
		exceptions.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if tokenID == uint32(k.ID) &&
		(k.Role == "manager" ||
			(k.Role == "coach" && k.ID == retrievedCoach.ID) ||
			(k.Role == "client" && k.ID == retrievedClient.ID)) {

		retrievedTask, err = task.GetTaskByID(server.Database, uint32(taskID))
		if err != nil {
			err = errors.New("500 Internal Server Error: cannot get task")
			exceptions.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(retrievedTask)

	} else {
		err = errors.New("Missing credentials: access restricted to Management, assigner coach and assignee client.")
		exceptions.ERROR(w, http.StatusUnauthorized, err)
		return
	}

}

func (server *Server) UpdateTask(w http.ResponseWriter, r *http.Request) {
	task := models.Task{}
	params := mux.Vars(r)
	taskID, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		err = errors.New("422 Unprocessable Entity Error: Cannot parse url")
		exceptions.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	retrievedTask, err := task.GetTaskByID(server.Database, uint32(taskID))

	u := models.User{}
	retrievedCoach, err := u.GetCoachByID(server.Database, retrievedTask.AssignerID)

	user := models.User{}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		err = errors.New("Invalid token: please login with valid credentials.")
		exceptions.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	k, err := user.GetUserByID(server.Database, tokenID)
	if err != nil {
		err = errors.New("500 Internal Server Error: cannot get user")
		exceptions.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if tokenID == uint32(k.ID) &&
		(k.Role == "manager" ||
			(k.Role == "coach" && k.ID == retrievedCoach.ID)) {

		taskToUpdate, err := task.GetTaskByID(server.Database, uint32(taskID))
		if err != nil {
			err = errors.New("500 Internal Server Error: cannot get task")
			exceptions.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			exceptions.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}

		err = json.Unmarshal(body, &taskToUpdate)
		if err != nil {
			err = errors.New("422 Unprocessable Entity Error: Cannot parse JSON data.")
			exceptions.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}

		t := models.Task{}

		alltasks, err := t.GetTasks(server.Database)

		for _, singleTask := range *alltasks {
			if singleTask.ID == taskToUpdate.ID {
				singleTask.Title = taskToUpdate.Title
				singleTask.Description = taskToUpdate.Description
				singleTask.Year = taskToUpdate.Year
				singleTask.Month = taskToUpdate.Month
				singleTask.Day = taskToUpdate.Day
				singleTask.Duration = taskToUpdate.Duration
				singleTask.DateTask = taskToUpdate.DateTask

				singleTaskMonth := ""
				if len(strconv.Itoa(task.Month)) == 1 {
					singleTaskMonth = "0" + strconv.Itoa(task.Month)
				} else {
					singleTaskMonth = strconv.Itoa(task.Month)
				}

				singleTaskDay := ""
				if len(strconv.Itoa(task.Day)) == 1 {
					singleTaskDay = "0" + strconv.Itoa(task.Day)
				} else {
					singleTaskDay = strconv.Itoa(task.Day)
				}

				t1, err := time.Parse(
					time.RFC3339,
					strconv.Itoa(task.Year)+"-"+singleTaskMonth+"-"+singleTaskDay+"T"+"00:00"+":00+00:00")

				if err != nil {
					err = errors.New("400 Bad Request Error: Cannot create task with this date")
					exceptions.ERROR(w, http.StatusBadRequest, err)
					return
				}

				if t1.Format(time.RFC3339) < time.Now().Format(time.RFC3339) {
					err = errors.New("400 Bad Request Error: Cannot create task with a past date")
					exceptions.ERROR(w, http.StatusBadRequest, err)
					return
				}

				singleTask.DateTask = strings.Replace(t1.Format(time.RFC3339)[:10], "T", " ", -1)

				server.Database.Save(singleTask)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(singleTask)
			}
		}
	} else {
		err = errors.New("Missing credentials: access restricted to Management and assigner coach.")
		exceptions.ERROR(w, http.StatusUnauthorized, err)
		return
	}
}

func (server *Server) DeleteTask(w http.ResponseWriter, r *http.Request) {

	task := models.Task{}
	params := mux.Vars(r)
	taskID, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		err = errors.New("422 Unprocessable Entity Error: Cannot parse url")
		exceptions.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	retrievedTask, err := task.GetTaskByID(server.Database, uint32(taskID))

	u := models.User{}
	retrievedCoach, err := u.GetCoachByID(server.Database, retrievedTask.AssignerID)

	user := models.User{}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		err = errors.New("Invalid token: please login with valid credentials.")
		exceptions.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	k, err := user.GetUserByID(server.Database, tokenID)
	if err != nil {
		err = errors.New("500 Internal Server Error: cannot get user")
		exceptions.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if tokenID == uint32(k.ID) &&
		(k.Role == "manager" ||
			(k.Role == "coach" && k.ID == retrievedCoach.ID)) {
		if _, err = task.DeleteTask(server.Database, uint32(taskID)); err != nil {
			err = errors.New("500 Internal Server Error: cannot delete task")
			exceptions.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(204)
		er := json.NewEncoder(w).Encode("")
		if er != nil {
			fmt.Fprintf(w, "%s", er.Error())
		}
	} else {
		err = errors.New("Missing credentials: access restricted to Management and assigner coach.")
		exceptions.ERROR(w, http.StatusUnauthorized, err)
		return
	}
}
