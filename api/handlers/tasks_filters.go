package handlers

import (
"encoding/json"
"errors"
"github.com/Alexandremerancienne/my_Sartorius/api/exceptions"
"github.com/Alexandremerancienne/my_Sartorius/api/models"
"net/http"
)

func (server *Server) ApplyTaskFilters(w http.ResponseWriter, tasks *[]models.Task, coach_id, client_id int, k models.User) {

	if coach_id != 0 && client_id == 0 {
		server.FilterTaskByAssignerID(w, tasks, coach_id, k)

	} else if coach_id == 0 && client_id != 0 {
		server.FilterTaskByAssigneeID(w, tasks, client_id, k)

	} else if coach_id != 0 && client_id != 0 {
		server.FilterTaskByCoachAndAssigneeID(w, tasks, coach_id, client_id, k)
	}
}

func (server *Server) FilterTaskByAssignerID(w http.ResponseWriter, tasks *[]models.Task, coach_id int, k models.User) {

	var output = []models.Task{}
	for _, task := range *tasks {
		if task.AssignerID == coach_id {
			output = append(output, task)
		}
	}

	if len(output) > 0 {

		if k.Role == "manager" {
			if len(output) > 0 {
				if len(output) == 1 {
					json.NewEncoder(w).Encode(output[0])
				} else {
					json.NewEncoder(w).Encode(output)
				}
			}
		} else if k.Role == "coach" {
			finalList := []models.Task{}
			for _, o := range output {
				if o.AssignerID == int(k.ID) {
					finalList = append(finalList, o)
				}
			}

			if len(finalList) > 0 {
				if len(finalList) == 1 {
					json.NewEncoder(w).Encode(finalList[0])
				} else {
					json.NewEncoder(w).Encode(finalList)
				}
			} else {
				err := errors.New("Missing Credentials: Please select another assigner ID as filter")
				exceptions.ERROR(w, http.StatusUnauthorized, err)
				return
			}

		} else if k.Role == "client" {
			finalList := []models.Task{}
			for _, o := range output {
				if o.AssigneeID == int(k.ID) {
					finalList = append(finalList, o)
				}
			}
			if len(finalList) > 0 {
				if len(finalList) == 1 {
					json.NewEncoder(w).Encode(finalList[0])
				} else {
					json.NewEncoder(w).Encode(finalList)
				}
			} else {
				err := errors.New("Missing Credentials: Please select another assigner ID as filter")
				exceptions.ERROR(w, http.StatusUnauthorized, err)
				return
			}
		}
	} else {
		err := errors.New("400 Error: no task found")
		exceptions.ERROR(w, http.StatusBadRequest, err)
		return
	}
}


func (server *Server) FilterTaskByAssigneeID(w http.ResponseWriter, tasks *[]models.Task, client_id int, k models.User) {

	var output = []models.Task{}
	for _, task := range *tasks {
		if task.AssigneeID == client_id {
			output = append(output, task)
		}
	}

	if len(output) > 0 {
		if k.Role == "manager" {
			if len(output) > 0 {
				if len(output) == 1 {
					json.NewEncoder(w).Encode(output[0])
				} else {
					json.NewEncoder(w).Encode(output)
				}
			}
		} else if k.Role == "coach" {
			finalList := []models.Task{}
			for _, o := range output {
				if o.AssignerID == int(k.ID) {
					finalList = append(finalList, o)
				}
			}

			if len(finalList) > 0 {
				if len(finalList) == 1 {
					json.NewEncoder(w).Encode(finalList[0])
				} else {
					json.NewEncoder(w).Encode(finalList)
				}
			} else {
				err := errors.New("Missing Credentials: Please select another assignee ID as filter")
				exceptions.ERROR(w, http.StatusUnauthorized, err)
				return
			}

		} else if k.Role == "client" {
			finalList := []models.Task{}

			for _, o := range output {
				if o.AssigneeID == int(k.ID) {
					finalList = append(finalList, o)
				}
			}
			if len(finalList) > 0 {
				if len(finalList) == 1 {
					json.NewEncoder(w).Encode(finalList[0])
				} else {
					json.NewEncoder(w).Encode(finalList)
				}
			} else {
				err := errors.New("Missing Credentials: Please select another assignee ID as filter")
				exceptions.ERROR(w, http.StatusUnauthorized, err)
				return
			}
		}
	} else {
		err := errors.New("400 Error: no task found")
		exceptions.ERROR(w, http.StatusBadRequest, err)
		return
	}
}

func (server *Server) FilterTaskByCoachAndAssigneeID(w http.ResponseWriter, tasks *[]models.Task, coach_id, client_id int, k models.User) {

	var output = []models.Task{}
	for _, task := range *tasks {
		if task.AssigneeID == client_id && task.AssignerID == coach_id {
			output = append(output, task)
		}
	}

	if len(output) > 0 {

		if k.Role == "manager" {
			if len(output) > 0 {
				if len(output) == 1 {
					json.NewEncoder(w).Encode(output[0])
				} else {
					json.NewEncoder(w).Encode(output)
				}
			}
		} else if k.Role == "coach" {
			finalList := []models.Task{}
			for _, o := range output {
				if o.AssignerID == int(k.ID) {
					finalList = append(finalList, o)
				}
			}

			if len(finalList) > 0 {
				if len(finalList) == 1 {
					json.NewEncoder(w).Encode(finalList[0])
				} else {
					json.NewEncoder(w).Encode(finalList)
				}
			} else {
				err := errors.New("401 Error: Access denied")
				exceptions.ERROR(w, http.StatusUnauthorized, err)
				return
			}

		} else if k.Role == "client" {
			finalList := []models.Task{}
			for _, o := range output {
				if o.AssigneeID == int(k.ID) {
					finalList = append(finalList, o)
				}
			}
			if len(finalList) > 0 {
				if len(finalList) == 1 {
					json.NewEncoder(w).Encode(finalList[0])
				} else {
					json.NewEncoder(w).Encode(finalList)
				}
			} else {
				err := errors.New("Missing Credentials: Please select another assigner ID or assignee ID as filter")
				exceptions.ERROR(w, http.StatusUnauthorized, err)
				return
			}
		}

	} else {
		err := errors.New("400 Error: no task found")
		exceptions.ERROR(w, http.StatusBadRequest, err)
		return
	}
}

