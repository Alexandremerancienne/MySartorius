package handlers

import (
	"encoding/json"
	"errors"
	"github.com/Alexandremerancienne/my_Sartorius/api/exceptions"
	"github.com/Alexandremerancienne/my_Sartorius/api/models"
	"net/http"
)

func (server *Server) ApplySessionFilters(w http.ResponseWriter, sessions *[]models.Session, coach_id, client_id int, k models.User) {

	if coach_id != 0 && client_id == 0 {
		server.FilterSessionByCoachID(w, sessions, coach_id, k)

	} else if coach_id == 0 && client_id != 0 {
		server.FilterSessionByClientID(w, sessions, client_id, k)

	} else if coach_id != 0 && client_id != 0 {
		server.FilterSessionByCoachAndClientID(w, sessions, coach_id, client_id, k)
	}
}

func (server *Server) FilterSessionByCoachID(w http.ResponseWriter, sessions *[]models.Session, coach_id int, k models.User) {

	var output = []models.Session{}
	for _, session := range *sessions {
		if session.CoachID == coach_id {
			output = append(output, session)
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
				finalList := []models.Session{}
			for _, o := range output {
				if o.CoachID == int(k.ID) {
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
				err := errors.New("Missing Credentials: Please select another coach ID as filter")
				exceptions.ERROR(w, http.StatusUnauthorized, err)
				return
			}

		} else if k.Role == "client" {
				finalList := []models.Session{}
				for _, o := range output {
				if o.ClientID == int(k.ID) {
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
				err := errors.New("Missing Credentials: Please select another coach ID as filter")
				exceptions.ERROR(w, http.StatusUnauthorized, err)
				return
			}
		}
	} else {
		err := errors.New("400 Error: no session found")
		exceptions.ERROR(w, http.StatusBadRequest, err)
		return
	}
}


func (server *Server) FilterSessionByClientID(w http.ResponseWriter, sessions *[]models.Session, client_id int, k models.User) {

	var output = []models.Session{}
	for _, session := range *sessions {
		if session.ClientID == client_id {
			output = append(output, session)
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
				finalList := []models.Session{}
			for _, o := range output {
				if o.CoachID == int(k.ID) {
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
				err := errors.New("Missing Credentials: Please select another client ID as filter")
				exceptions.ERROR(w, http.StatusUnauthorized, err)
				return
			}

		} else if k.Role == "client" {
				finalList := []models.Session{}

			for _, o := range output {
				if o.ClientID == int(k.ID) {
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
				err := errors.New("Missing Credentials: Please select another client ID as filter")
				exceptions.ERROR(w, http.StatusUnauthorized, err)
				return
			}
		}
	} else {
		err := errors.New("400 Error: no session found")
		exceptions.ERROR(w, http.StatusBadRequest, err)
		return
	}
}

func (server *Server) FilterSessionByCoachAndClientID(w http.ResponseWriter, sessions *[]models.Session, coach_id, client_id int, k models.User) {

	var output = []models.Session{}
	for _, session := range *sessions {
		if session.ClientID == client_id && session.CoachID == coach_id {
			output = append(output, session)
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
				finalList := []models.Session{}
			for _, o := range output {
				if o.CoachID == int(k.ID) {
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
				finalList := []models.Session{}
			for _, o := range output {
				if o.ClientID == int(k.ID) {
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
				err := errors.New("Missing Credentials: Please select another coach ID or client ID as filter")
				exceptions.ERROR(w, http.StatusUnauthorized, err)
				return
			}
		}

	} else {
		err := errors.New("400 Error: no session found")
		exceptions.ERROR(w, http.StatusBadRequest, err)
		return
	}
}
