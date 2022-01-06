package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Alexandremerancienne/my_Sartorius/api/auth"
	"github.com/Alexandremerancienne/my_Sartorius/api/exceptions"
	"github.com/Alexandremerancienne/my_Sartorius/api/models"
	"github.com/gorilla/mux"
)

func (server *Server) CreateReminder(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	taskID, err := strconv.ParseUint(params["task_id"], 10, 32)
	if err != nil {
		err = errors.New("422 Unprocessable Entity Error: Cannot parse url")
		exceptions.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

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

	task := models.Task{}
	retrievedTask, err := task.GetTaskByID(server.Database, uint32(taskID))
	if err != nil {
		err = errors.New("500 Internal Server Error: cannot get task")
		exceptions.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	u := models.User{}
	reminderOwner, err := u.GetClientByID(server.Database, retrievedTask.AssigneeID)
	if err != nil {
		err = errors.New("500 Internal Server Error: cannot get client")
		exceptions.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if tokenID == uint32(k.ID) && k.Role == "client" && k.ID == reminderOwner.ID {

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			exceptions.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}

		reminder := models.Reminder{}
		if err := json.Unmarshal(body, &reminder); err != nil {
			err = errors.New("422 Unprocessable Entity Error: Cannot parse JSON data.")
			exceptions.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}

		newReminder, err := reminder.CreateReminder(server.Database, taskID)
		if err != nil {
			err = errors.New("500 Internal Server Error: cannot create reminder")
			exceptions.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(newReminder)

	} else {
		err = errors.New("Missing credentials: You have not access to this task and its reminders.")
		exceptions.ERROR(w, http.StatusUnauthorized, err)
		return
	}
}

func (server *Server) GetReminders(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	taskID, err := strconv.ParseUint(params["task_id"], 10, 32)
	if err != nil {
		err = errors.New("422 Unprocessable Entity Error: Cannot parse url")
		exceptions.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	reminders := models.Reminder{}
	retrievedReminders, err := reminders.GetReminders(server.Database, taskID)
	if err != nil {
		err = errors.New("500 Internal Server Error: cannot get reminders")
		exceptions.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	remindersSlice := []models.Reminder{}

	if len(*retrievedReminders) == 0 {

		task := models.Task{}
		retrievedTask, err := task.GetTaskByID(server.Database, uint32(taskID))
		if err != nil {
			err = errors.New("500 Internal Server Error: cannot get task")
			exceptions.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		u := models.User{}
		reminderOwner, err := u.GetClientByID(server.Database, retrievedTask.AssigneeID)
		if err != nil {
			err = errors.New("500 Internal Server Error: cannot get client")
			exceptions.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		tokenID, err := auth.ExtractTokenID(r)
		if err != nil {
			err = errors.New("Invalid token: please login with valid credentials.")
			exceptions.ERROR(w, http.StatusUnauthorized, err)
			return
		}

		h := models.User{}
		k, err := h.GetUserByID(server.Database, tokenID)
		if err != nil {
			err = errors.New("500 Internal Server Error: cannot get user")
			exceptions.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		if tokenID == uint32(k.ID) && k.Role == "client" && k.ID == reminderOwner.ID {
			err = errors.New("No reminders for this task.")
			exceptions.ERROR(w, http.StatusOK, err)
			return
		} else {
			err = errors.New("Missing credentials: You have not access to this task and its reminders.")
			exceptions.ERROR(w, http.StatusUnauthorized, err)
			return
		}

	} else {

		for _, reminder := range *retrievedReminders {

			task := models.Task{}
			retrievedTask, err := task.GetTaskByID(server.Database, uint32(taskID))
			if err != nil {
				err = errors.New("500 Internal Server Error: cannot get task")
				exceptions.ERROR(w, http.StatusInternalServerError, err)
				return
			}

			u := models.User{}
			reminderOwner, err := u.GetClientByID(server.Database, retrievedTask.AssigneeID)

			if err != nil {
				err = errors.New("500 Internal Server Error: cannot get client")
				exceptions.ERROR(w, http.StatusInternalServerError, err)
				return
			}

			tokenID, err := auth.ExtractTokenID(r)
			if err != nil {
				err = errors.New("Invalid token: please login with valid credentials.")
				exceptions.ERROR(w, http.StatusUnauthorized, err)
				return
			}

			l := models.User{}

			k, err := l.GetUserByID(server.Database, tokenID)
			if err != nil {
				err = errors.New("500 Internal Server Error: cannot get user")
				exceptions.ERROR(w, http.StatusInternalServerError, err)
				return
			}

			if tokenID == uint32(k.ID) && k.Role == "client" && k.ID == reminderOwner.ID {
				remindersSlice = append(remindersSlice, reminder)
			} else {
				err = errors.New("Missing credentials: You have not access to this task and its reminders.")
				exceptions.ERROR(w, http.StatusUnauthorized, err)
				return
			}

		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(remindersSlice)
	}
}

func (server *Server) GetReminder(w http.ResponseWriter, r *http.Request) {

	reminder := models.Reminder{}
	params := mux.Vars(r)
	taskID, err := strconv.ParseUint(params["task_id"], 10, 32)
	if err != nil {
		err = errors.New("422 Unprocessable Entity Error: Cannot parse url")
		exceptions.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	reminderId, err := strconv.ParseUint(params["reminder_id"], 10, 32)
	if err != nil {
		err = errors.New("422 Unprocessable Entity Error: Cannot parse url")
		exceptions.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

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

	task := models.Task{}
	retrievedTask, err := task.GetTaskByID(server.Database, uint32(taskID))
	if err != nil {
		err = errors.New("500 Internal Server Error: cannot get task")
		exceptions.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	u := models.User{}
	reminderOwner, err := u.GetClientByID(server.Database, retrievedTask.AssigneeID)
	if err != nil {
		err = errors.New("500 Internal Server Error: cannot get client")
		exceptions.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if tokenID == uint32(k.ID) && k.Role == "client" && k.ID == reminderOwner.ID {
		retrievedReminder, err := reminder.GetReminderByID(server.Database, uint32(reminderId), uint32(taskID))
		if err != nil {
			err = errors.New("500 Internal Server Error: cannot get reminder")
			exceptions.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(retrievedReminder)
	} else {
		err = errors.New("Missing credentials: You have not access to this task and its reminders.")
		exceptions.ERROR(w, http.StatusUnauthorized, err)
		return
	}
}

func (server *Server) UpdateReminder(w http.ResponseWriter, r *http.Request) {

	reminder := models.Reminder{}
	params := mux.Vars(r)
	taskID, err := strconv.ParseUint(params["task_id"], 10, 32)
	if err != nil {
		err = errors.New("422 Unprocessable Entity Error: Cannot parse url")
		exceptions.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	reminderId, err := strconv.ParseUint(params["reminder_id"], 10, 32)
	if err != nil {
		err = errors.New("422 Unprocessable Entity Error: Cannot parse url")
		exceptions.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

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

	task := models.Task{}
	retrievedTask, err := task.GetTaskByID(server.Database, uint32(taskID))
	if err != nil {
		err = errors.New("500 Internal Server Error: cannot get task")
		exceptions.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	u := models.User{}
	reminderOwner, err := u.GetClientByID(server.Database, retrievedTask.AssigneeID)
	if err != nil {
		err = errors.New("500 Internal Server Error: cannot get client")
		exceptions.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if tokenID == uint32(k.ID) && k.Role == "client" && k.ID == reminderOwner.ID {
		w.Header().Set("Content-Type", "application/json")

		reminderToUpdate, err := reminder.GetReminderByID(server.Database, uint32(reminderId), uint32(taskID))
		if err != nil {
			err = errors.New("500 Internal Server Error: cannot get reminder")
			exceptions.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			exceptions.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}

		err = json.Unmarshal(body, &reminderToUpdate)
		if err != nil {
			err = errors.New("422 Unprocessable Entity Error: Cannot parse JSON data.")
			exceptions.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}

		r := models.Reminder{}

		allreminders, err := r.GetReminders(server.Database, taskID)

		for _, singleReminder := range *allreminders {
			if singleReminder.ID == reminderToUpdate.ID {
				singleReminder.Description = reminderToUpdate.Description
				server.Database.Save(singleReminder)
				json.NewEncoder(w).Encode(singleReminder)
			}
		}
	} else {
		err = errors.New("Missing credentials: You have not access to this task and its reminders.")
		exceptions.ERROR(w, http.StatusUnauthorized, err)
		return
	}
}

func (server *Server) DeleteReminder(w http.ResponseWriter, r *http.Request) {

	reminder := models.Reminder{}
	params := mux.Vars(r)
	taskID, err := strconv.ParseUint(params["task_id"], 10, 32)
	if err != nil {
		err = errors.New("422 Unprocessable Entity Error: Cannot parse url")
		exceptions.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	reminderId, err := strconv.ParseUint(params["reminder_id"], 10, 32)
	if err != nil {
		err = errors.New("422 Unprocessable Entity Error: Cannot parse url")
		exceptions.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

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

	task := models.Task{}
	retrievedTask, err := task.GetTaskByID(server.Database, uint32(taskID))
	if err != nil {
		err = errors.New("500 Internal Server Error: cannot get task")
		exceptions.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	u := models.User{}
	reminderOwner, err := u.GetClientByID(server.Database, retrievedTask.AssigneeID)
	if err != nil {
		err = errors.New("500 Internal Server Error: cannot get client")
		exceptions.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if tokenID == uint32(k.ID) && k.Role == "client" && k.ID == reminderOwner.ID {
		if _, err = reminder.DeleteReminder(server.Database, uint32(reminderId), uint32(taskID)); err != nil {
			err = errors.New("500 Internal Server Error: cannot delete reminder")
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
		err = errors.New("Missing credentials: You have not access to this task and its reminders.")
		exceptions.ERROR(w, http.StatusUnauthorized, err)
		return
	}
}

