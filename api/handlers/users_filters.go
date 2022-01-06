package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/Alexandremerancienne/my_Sartorius/api/exceptions"
	"github.com/Alexandremerancienne/my_Sartorius/api/models"
)

func (server *Server) ApplyCoachFilters(w http.ResponseWriter, coaches *[]models.User, user models.User, id int, last_name, first_name, last_name_contains, first_name_contains string) {
	if id == 0 && last_name == "" && first_name == "" && last_name_contains == "" && first_name_contains == "" {
		coachesWithoutPassword := []UserWithoutPassword{}
		for _, c := range *coaches {
			newUserWithoutPassword := new(UserWithoutPassword)
			newUserWithoutPassword.ID = int(c.ID)
			newUserWithoutPassword.FirstName = c.FirstName
			newUserWithoutPassword.LastName = c.LastName
			newUserWithoutPassword.Email = c.Email
			newUserWithoutPassword.PhoneNumber = c.PhoneNumber
			newUserWithoutPassword.Role = c.Role
			coachesWithoutPassword = append(coachesWithoutPassword, *newUserWithoutPassword)
		}
		json.NewEncoder(w).Encode(coachesWithoutPassword)

	} else if id != 0 && last_name == "" && first_name == "" && last_name_contains == "" && first_name_contains == "" {
		server.FilterCoachByID(user, id, w)

	} else if id != 0 && last_name != "" && first_name == "" && last_name_contains == "" && first_name_contains == "" {
		server.FilterCoachByIDAndLastName(user, id, w, last_name)

	} else if id != 0 && last_name == "" && first_name != "" && last_name_contains == "" && first_name_contains == "" {
		server.FilterCoachByIDAndFirstName(user, id, w, first_name)

	} else if id != 0 && last_name != "" && first_name != "" && last_name_contains == "" && first_name_contains == "" {
		server.FilterCoachByIDFirstNameAndLastName(user, id, w, first_name, last_name)

	} else if id == 0 && last_name != "" && first_name == "" && last_name_contains == "" && first_name_contains == "" {
		server.FilterCoachByLastName(coaches, w, last_name)

	} else if id == 0 && last_name == "" && first_name != "" && last_name_contains == "" && first_name_contains == "" {
		server.FilterCoachByFirstName(coaches, w, first_name)

	} else if id == 0 && last_name != "" && first_name != "" && last_name_contains == "" && first_name_contains == "" {
		server.FilterCoachByLastNameAndFirstName(coaches, w, first_name, last_name)

	} else if id != 0 && last_name == "" && first_name == "" && last_name_contains != "" && first_name_contains == "" {
		server.FilterCoachByIDWithLastNameContains(user, id, w, last_name_contains)

	} else if id != 0 && last_name == "" && first_name == "" && last_name_contains == "" && first_name_contains != "" {
		server.FilterCoachByIDWithFirstNameContains(user, id, w, first_name_contains)

	} else if id != 0 && last_name == "" && first_name == "" && last_name_contains != "" && first_name_contains != "" {
		server.FilterCoachByIDWithFirstNameContainsAndLastNameContains(user, id, w, last_name_contains, first_name_contains)

	} else if id == 0 && last_name == "" && first_name == "" && last_name_contains != "" && first_name_contains == "" {
		server.FilterCoachWithLastNameContains(coaches, w, last_name_contains)

	} else if id == 0 && last_name == "" && first_name == "" && last_name_contains == "" && first_name_contains != "" {
		server.FilterCoachWithFirstNameContains(coaches, w, first_name_contains)

	} else if id == 0 && last_name == "" && first_name == "" && last_name_contains != "" && first_name_contains != "" {
		server.FilterCoachWithLastNameContainsAndFirstNameContains(coaches, w, last_name_contains, first_name_contains)
	}
}

func (server *Server) ApplyClientFilters(w http.ResponseWriter, clients *[]models.User, user models.User, id int, last_name, first_name, last_name_contains, first_name_contains string) {
	if id == 0 && last_name == "" && first_name == "" && last_name_contains == "" && first_name_contains == "" {
		clientsWithoutPassword := []UserWithoutPassword{}
		for _, c := range *clients {
			newUserWithoutPassword := new(UserWithoutPassword)
			newUserWithoutPassword.ID = int(c.ID)
			newUserWithoutPassword.FirstName = c.FirstName
			newUserWithoutPassword.LastName = c.LastName
			newUserWithoutPassword.Email = c.Email
			newUserWithoutPassword.PhoneNumber = c.PhoneNumber
			newUserWithoutPassword.Role = c.Role
			clientsWithoutPassword = append(clientsWithoutPassword, *newUserWithoutPassword)
		}
		json.NewEncoder(w).Encode(clientsWithoutPassword)

	} else if id != 0 && last_name == "" && first_name == "" && last_name_contains == "" && first_name_contains == "" {
		server.FilterClientByID(user, id, w)

	} else if id != 0 && last_name != "" && first_name == "" && last_name_contains == "" && first_name_contains == "" {
		server.FilterClientByIDAndLastName(user, id, w, last_name)

	} else if id != 0 && last_name == "" && first_name != "" && last_name_contains == "" && first_name_contains == "" {
		server.FilterClientByIDAndFirstName(user, id, w, first_name)

	} else if id != 0 && last_name != "" && first_name != "" && last_name_contains == "" && first_name_contains == "" {
		server.FilterClientByIDFirstNameAndLastName(user, id, w, first_name, last_name)

	} else if id == 0 && last_name != "" && first_name == "" && last_name_contains == "" && first_name_contains == "" {
		server.FilterClientByLastName(clients, w, last_name)

	} else if id == 0 && last_name == "" && first_name != "" && last_name_contains == "" && first_name_contains == "" {
		server.FilterClientByFirstName(clients, w, first_name)

	} else if id == 0 && last_name != "" && first_name != "" && last_name_contains == "" && first_name_contains == "" {
		server.FilterClientByLastNameAndFirstName(clients, w, first_name, last_name)

	} else if id != 0 && last_name == "" && first_name == "" && last_name_contains != "" && first_name_contains == "" {
		server.FilterClientByIDWithLastNameContains(user, id, w, last_name_contains)

	} else if id != 0 && last_name == "" && first_name == "" && last_name_contains == "" && first_name_contains != "" {
		server.FilterClientByIDWithFirstNameContains(user, id, w, first_name_contains)

	} else if id != 0 && last_name == "" && first_name == "" && last_name_contains != "" && first_name_contains != "" {
		server.FilterClientByIDWithFirstNameContainsAndLastNameContains(user, id, w, last_name_contains, first_name_contains)

	} else if id == 0 && last_name == "" && first_name == "" && last_name_contains != "" && first_name_contains == "" {
		server.FilterClientWithLastNameContains(clients, w, last_name_contains)

	} else if id == 0 && last_name == "" && first_name == "" && last_name_contains == "" && first_name_contains != "" {
		server.FilterClientWithFirstNameContains(clients, w, first_name_contains)

	} else if id == 0 && last_name == "" && first_name == "" && last_name_contains != "" && first_name_contains != "" {
		server.FilterClientWithLastNameContainsAndFirstNameContains(clients, w, last_name_contains, first_name_contains)
	}
}

func (server *Server) FilterCoachByID(u models.User, id int, w http.ResponseWriter) {
	coach, err := u.GetCoachByID(server.Database, id)
	if err != nil {
		err = errors.New("404 Error: coach not found")
		exceptions.ERROR(w, http.StatusNotFound, err)
		return
	}
	newUserWithoutPassword := new(UserWithoutPassword)
	newUserWithoutPassword.ID = int(coach.ID)
	newUserWithoutPassword.FirstName = coach.FirstName
	newUserWithoutPassword.LastName = coach.LastName
	newUserWithoutPassword.Email = coach.Email
	newUserWithoutPassword.PhoneNumber = coach.PhoneNumber
	newUserWithoutPassword.Role = coach.Role
	json.NewEncoder(w).Encode(newUserWithoutPassword)
}

func (server *Server) FilterCoachByLastName(u *[]models.User, w http.ResponseWriter, lastName string) {
	var output = []models.User{}
	for _, c := range *u {
		if strings.ToLower(c.LastName) == strings.ToLower(lastName) {
			output = append(output, c)
		}
	}
	if len(output) > 0 {

		coachesWithoutPassword := []UserWithoutPassword{}

		for _, coach := range output {
			newUserWithoutPassword := new(UserWithoutPassword)
			newUserWithoutPassword.ID = int(coach.ID)
			newUserWithoutPassword.FirstName = coach.FirstName
			newUserWithoutPassword.LastName = coach.LastName
			newUserWithoutPassword.Email = coach.Email
			newUserWithoutPassword.PhoneNumber = coach.PhoneNumber
			newUserWithoutPassword.Role = coach.Role
			coachesWithoutPassword = append(coachesWithoutPassword, *newUserWithoutPassword)
		}

		if len(output) == 1 {
			json.NewEncoder(w).Encode(coachesWithoutPassword[0])
		} else {
			json.NewEncoder(w).Encode(coachesWithoutPassword)
		}
	} else {
		err := errors.New("404 Error: no coach found")
		exceptions.ERROR(w, http.StatusNotFound, err)
		return
	}
}

func (server *Server) FilterCoachByFirstName(u *[]models.User, w http.ResponseWriter, firstName string) {
	var output = []models.User{}
	for _, c := range *u {
		if strings.ToLower(c.FirstName) == strings.ToLower(firstName) {
			output = append(output, c)
		}
	}
	if len(output) > 0 {

		coachesWithoutPassword := []UserWithoutPassword{}

		for _, coach := range output {
			newUserWithoutPassword := new(UserWithoutPassword)
			newUserWithoutPassword.ID = int(coach.ID)
			newUserWithoutPassword.FirstName = coach.FirstName
			newUserWithoutPassword.LastName = coach.LastName
			newUserWithoutPassword.Email = coach.Email
			newUserWithoutPassword.PhoneNumber = coach.PhoneNumber
			newUserWithoutPassword.Role = coach.Role
			coachesWithoutPassword = append(coachesWithoutPassword, *newUserWithoutPassword)
		}
		if len(output) == 1 {
			json.NewEncoder(w).Encode(coachesWithoutPassword[0])
		} else {
			json.NewEncoder(w).Encode(coachesWithoutPassword)
		}
	} else {
		err := errors.New("404 Error: no coach found")
		exceptions.ERROR(w, http.StatusNotFound, err)
		return
	}
}

func (server *Server) FilterCoachByLastNameAndFirstName(u *[]models.User, w http.ResponseWriter, firstName, lastName string) {
	var output = []models.User{}
	for _, c := range *u {
		if strings.ToLower(c.FirstName) == strings.ToLower(firstName) {
			output = append(output, c)
		}
	}
	var finalList = []models.User{}
	if len(output) > 0 {
		for _, o := range output {
			if strings.ToLower(o.LastName) == strings.ToLower(lastName) {
				finalList = append(finalList, o)
			}
		}
		if len(finalList) > 0 {

			coachesWithoutPassword := []UserWithoutPassword{}
			for _, coach := range finalList {
				newUserWithoutPassword := new(UserWithoutPassword)
				newUserWithoutPassword.ID = int(coach.ID)
				newUserWithoutPassword.FirstName = coach.FirstName
				newUserWithoutPassword.LastName = coach.LastName
				newUserWithoutPassword.Email = coach.Email
				newUserWithoutPassword.PhoneNumber = coach.PhoneNumber
				newUserWithoutPassword.Role = coach.Role
				coachesWithoutPassword = append(coachesWithoutPassword, *newUserWithoutPassword)
			}

			if len(finalList) == 1 {
				json.NewEncoder(w).Encode(coachesWithoutPassword[0])
			} else {
				json.NewEncoder(w).Encode(coachesWithoutPassword)
			}
		} else {
			err := errors.New("404 Error: no coach found")
			exceptions.ERROR(w, http.StatusNotFound, err)
			return
		}
	}
}

func (server *Server) FilterCoachByIDAndLastName(u models.User, id int, w http.ResponseWriter, lastName string) {
	coach, err := u.GetCoachByID(server.Database, id)
	if err != nil {
		err = errors.New("404 Error: coach not found")
		exceptions.ERROR(w, http.StatusNotFound, err)
		return
	}

	if strings.ToLower(coach.LastName) == strings.ToLower(lastName) {

		newUserWithoutPassword := new(UserWithoutPassword)
		newUserWithoutPassword.ID = int(coach.ID)
		newUserWithoutPassword.FirstName = coach.FirstName
		newUserWithoutPassword.LastName = coach.LastName
		newUserWithoutPassword.Email = coach.Email
		newUserWithoutPassword.PhoneNumber = coach.PhoneNumber
		newUserWithoutPassword.Role = coach.Role

		json.NewEncoder(w).Encode(newUserWithoutPassword)
	} else {
		err = errors.New("404 Error: coach not found")
		exceptions.ERROR(w, http.StatusNotFound, err)
		return
	}
}

func (server *Server) FilterCoachByIDAndFirstName(u models.User, id int, w http.ResponseWriter, firstName string) {
	coach, err := u.GetCoachByID(server.Database, id)
	if err != nil {
		err = errors.New("404 Error: coach not found")
		exceptions.ERROR(w, http.StatusNotFound, err)
		return
	}

	if strings.ToLower(coach.FirstName) == strings.ToLower(firstName) {

		newUserWithoutPassword := new(UserWithoutPassword)
		newUserWithoutPassword.ID = int(coach.ID)
		newUserWithoutPassword.FirstName = coach.FirstName
		newUserWithoutPassword.LastName = coach.LastName
		newUserWithoutPassword.Email = coach.Email
		newUserWithoutPassword.PhoneNumber = coach.PhoneNumber
		newUserWithoutPassword.Role = coach.Role

		json.NewEncoder(w).Encode(newUserWithoutPassword)
	} else {
		err = errors.New("404 Error: coach not found")
		exceptions.ERROR(w, http.StatusNotFound, err)
		return
	}
}

func (server *Server) FilterCoachByIDFirstNameAndLastName(u models.User, id int, w http.ResponseWriter, firstName, lastName string) {
	coach, err := u.GetCoachByID(server.Database, id)
	if err != nil {
		err = errors.New("404 Error: coach not found")
		exceptions.ERROR(w, http.StatusNotFound, err)
		return
	}

	if strings.ToLower(coach.FirstName) == strings.ToLower(firstName) && strings.ToLower(coach.LastName) == strings.ToLower(lastName) {

		newUserWithoutPassword := new(UserWithoutPassword)
		newUserWithoutPassword.ID = int(coach.ID)
		newUserWithoutPassword.FirstName = coach.FirstName
		newUserWithoutPassword.LastName = coach.LastName
		newUserWithoutPassword.Email = coach.Email
		newUserWithoutPassword.PhoneNumber = coach.PhoneNumber
		newUserWithoutPassword.Role = coach.Role

		json.NewEncoder(w).Encode(newUserWithoutPassword)
	} else {
		err = errors.New("404 Error: coach not found")
		exceptions.ERROR(w, http.StatusNotFound, err)
		return
	}
}

func (server *Server) FilterCoachByIDWithLastNameContains(u models.User, id int, w http.ResponseWriter, c string) {
	coach, err := u.GetCoachByID(server.Database, id)
	if err != nil {
		err = errors.New("404 Error: coach not found")
		exceptions.ERROR(w, http.StatusNotFound, err)
		return
	}
	if strings.Contains(strings.ToLower(coach.LastName), strings.ToLower(c)) {

		newUserWithoutPassword := new(UserWithoutPassword)
		newUserWithoutPassword.ID = int(coach.ID)
		newUserWithoutPassword.FirstName = coach.FirstName
		newUserWithoutPassword.LastName = coach.LastName
		newUserWithoutPassword.Email = coach.Email
		newUserWithoutPassword.PhoneNumber = coach.PhoneNumber
		newUserWithoutPassword.Role = coach.Role

		json.NewEncoder(w).Encode(newUserWithoutPassword)
	} else {
		err = errors.New("404 Error: coach not found")
		exceptions.ERROR(w, http.StatusNotFound, err)
		return
	}
}

func (server *Server) FilterCoachByIDWithFirstNameContains(u models.User, id int, w http.ResponseWriter, c string) {
	coach, err := u.GetCoachByID(server.Database, id)
	if err != nil {
		err = errors.New("404 Error: coach not found")
		exceptions.ERROR(w, http.StatusNotFound, err)
		return
	}
	if strings.Contains(strings.ToLower(coach.FirstName), strings.ToLower(c)) {

		newUserWithoutPassword := new(UserWithoutPassword)
		newUserWithoutPassword.ID = int(coach.ID)
		newUserWithoutPassword.FirstName = coach.FirstName
		newUserWithoutPassword.LastName = coach.LastName
		newUserWithoutPassword.Email = coach.Email
		newUserWithoutPassword.PhoneNumber = coach.PhoneNumber
		newUserWithoutPassword.Role = coach.Role

		json.NewEncoder(w).Encode(newUserWithoutPassword)
	} else {
		err = errors.New("404 Error: coach not found")
		exceptions.ERROR(w, http.StatusNotFound, err)
		return
	}
}

func (server *Server) FilterCoachByIDWithFirstNameContainsAndLastNameContains(u models.User, id int, w http.ResponseWriter, lastName, firstName string) {
	coach, err := u.GetCoachByID(server.Database, id)
	if err != nil {
		err = errors.New("404 Error: coach not found")
		exceptions.ERROR(w, http.StatusNotFound, err)
		return
	}
	if strings.Contains(strings.ToLower(coach.FirstName), strings.ToLower(firstName)) && strings.Contains(strings.ToLower(coach.LastName), strings.ToLower(lastName)) {

		newUserWithoutPassword := new(UserWithoutPassword)
		newUserWithoutPassword.ID = int(coach.ID)
		newUserWithoutPassword.FirstName = coach.FirstName
		newUserWithoutPassword.LastName = coach.LastName
		newUserWithoutPassword.Email = coach.Email
		newUserWithoutPassword.PhoneNumber = coach.PhoneNumber
		newUserWithoutPassword.Role = coach.Role

		json.NewEncoder(w).Encode(newUserWithoutPassword)
	} else {
		err = errors.New("404 Error: coach not found")
		exceptions.ERROR(w, http.StatusNotFound, err)
		return
	}
}

func (server *Server) FilterCoachWithLastNameContains(u *[]models.User, w http.ResponseWriter, lastName string) {
	var output = []models.User{}
	for _, c := range *u {
		if strings.Contains(strings.ToLower(c.LastName), strings.ToLower(lastName)) {
			output = append(output, c)
		}
	}
	if len(output) > 0 {

		coachesWithoutPassword := []UserWithoutPassword{}
		for _, coach := range output {
			newUserWithoutPassword := new(UserWithoutPassword)
			newUserWithoutPassword.ID = int(coach.ID)
			newUserWithoutPassword.FirstName = coach.FirstName
			newUserWithoutPassword.LastName = coach.LastName
			newUserWithoutPassword.Email = coach.Email
			newUserWithoutPassword.PhoneNumber = coach.PhoneNumber
			newUserWithoutPassword.Role = coach.Role
			coachesWithoutPassword = append(coachesWithoutPassword, *newUserWithoutPassword)
		}

		if len(output) == 1 {
			json.NewEncoder(w).Encode(coachesWithoutPassword[0])
		} else {
			json.NewEncoder(w).Encode(coachesWithoutPassword)
		}
	} else {
		err := errors.New("404 Error: no coach found")
		exceptions.ERROR(w, http.StatusNotFound, err)
		return
	}
}

func (server *Server) FilterCoachWithFirstNameContains(u *[]models.User, w http.ResponseWriter, firstName string) {
	var output = []models.User{}
	for _, c := range *u {
		if strings.Contains(strings.ToLower(c.FirstName), strings.ToLower(firstName)) {
			output = append(output, c)
		}
	}
	if len(output) > 0 {

		coachesWithoutPassword := []UserWithoutPassword{}
		for _, coach := range output {
			newUserWithoutPassword := new(UserWithoutPassword)
			newUserWithoutPassword.ID = int(coach.ID)
			newUserWithoutPassword.FirstName = coach.FirstName
			newUserWithoutPassword.LastName = coach.LastName
			newUserWithoutPassword.Email = coach.Email
			newUserWithoutPassword.PhoneNumber = coach.PhoneNumber
			newUserWithoutPassword.Role = coach.Role
			coachesWithoutPassword = append(coachesWithoutPassword, *newUserWithoutPassword)
		}

		if len(output) == 1 {
			json.NewEncoder(w).Encode(coachesWithoutPassword[0])
		} else {
			json.NewEncoder(w).Encode(coachesWithoutPassword)
		}
	} else {
		err := errors.New("404 Error: no coach found")
		exceptions.ERROR(w, http.StatusNotFound, err)
		return
	}
}

func (server *Server) FilterCoachWithLastNameContainsAndFirstNameContains(u *[]models.User, w http.ResponseWriter, lastName, firstName string) {
	var output = []models.User{}
	for _, c := range *u {
		if strings.Contains(strings.ToLower(c.FirstName), strings.ToLower(firstName)) {
			output = append(output, c)
		}
	}
	var finalList = []models.User{}
	if len(output) > 0 {
		for _, o := range output {
			if strings.Contains(strings.ToLower(o.LastName), strings.ToLower(lastName)) {
				finalList = append(finalList, o)
			}
		}
		if len(finalList) > 0 {

			coachesWithoutPassword := []UserWithoutPassword{}
			for _, coach := range finalList {
				newUserWithoutPassword := new(UserWithoutPassword)
				newUserWithoutPassword.ID = int(coach.ID)
				newUserWithoutPassword.FirstName = coach.FirstName
				newUserWithoutPassword.LastName = coach.LastName
				newUserWithoutPassword.Email = coach.Email
				newUserWithoutPassword.PhoneNumber = coach.PhoneNumber
				newUserWithoutPassword.Role = coach.Role
				coachesWithoutPassword = append(coachesWithoutPassword, *newUserWithoutPassword)
			}

			if len(output) == 1 {
				json.NewEncoder(w).Encode(coachesWithoutPassword[0])
			} else {
				json.NewEncoder(w).Encode(coachesWithoutPassword)
			}
		} else {
			err := errors.New("404 Error: no coach found")
			exceptions.ERROR(w, http.StatusNotFound, err)
			return
		}
	}
}

func (server *Server) FilterClientByID(u models.User, id int, w http.ResponseWriter) {
	client, err := u.GetClientByID(server.Database, id)
	if err != nil {
		err = errors.New("404 Error: client not found")
		exceptions.ERROR(w, http.StatusNotFound, err)
		return
	}

	newUserWithoutPassword := new(UserWithoutPassword)
	newUserWithoutPassword.ID = int(client.ID)
	newUserWithoutPassword.FirstName = client.FirstName
	newUserWithoutPassword.LastName = client.LastName
	newUserWithoutPassword.Email = client.Email
	newUserWithoutPassword.PhoneNumber = client.PhoneNumber
	newUserWithoutPassword.Role = client.Role

	json.NewEncoder(w).Encode(newUserWithoutPassword)
}

func (server *Server) FilterClientByLastName(u *[]models.User, w http.ResponseWriter, lastName string) {
	var output = []models.User{}
	for _, c := range *u {
		if strings.ToLower(c.LastName) == strings.ToLower(lastName) {
			output = append(output, c)
		}
	}
	if len(output) > 0 {

		clientsWithoutPassword := []UserWithoutPassword{}

		for _, client := range output {
			newUserWithoutPassword := new(UserWithoutPassword)
			newUserWithoutPassword.ID = int(client.ID)
			newUserWithoutPassword.FirstName = client.FirstName
			newUserWithoutPassword.LastName = client.LastName
			newUserWithoutPassword.Email = client.Email
			newUserWithoutPassword.PhoneNumber = client.PhoneNumber
			newUserWithoutPassword.Role = client.Role
			clientsWithoutPassword = append(clientsWithoutPassword, *newUserWithoutPassword)
		}

		if len(output) == 1 {
			json.NewEncoder(w).Encode(clientsWithoutPassword[0])
		} else {
			json.NewEncoder(w).Encode(clientsWithoutPassword)
		}
	} else {
		err := errors.New("404 Error: no client found")
		exceptions.ERROR(w, http.StatusNotFound, err)
		return
	}
}

func (server *Server) FilterClientByFirstName(u *[]models.User, w http.ResponseWriter, firstName string) {
	var output = []models.User{}
	for _, c := range *u {
		if strings.ToLower(c.FirstName) == strings.ToLower(firstName) {
			output = append(output, c)
		}
	}
	if len(output) > 0 {

		clientsWithoutPassword := []UserWithoutPassword{}

		for _, client := range output {
			newUserWithoutPassword := new(UserWithoutPassword)
			newUserWithoutPassword.ID = int(client.ID)
			newUserWithoutPassword.FirstName = client.FirstName
			newUserWithoutPassword.LastName = client.LastName
			newUserWithoutPassword.Email = client.Email
			newUserWithoutPassword.PhoneNumber = client.PhoneNumber
			newUserWithoutPassword.Role = client.Role
			clientsWithoutPassword = append(clientsWithoutPassword, *newUserWithoutPassword)
		}

		if len(output) == 1 {
			json.NewEncoder(w).Encode(clientsWithoutPassword[0])
		} else {
			json.NewEncoder(w).Encode(clientsWithoutPassword)
		}
	} else {
		err := errors.New("404 Error: no client found")
		exceptions.ERROR(w, http.StatusNotFound, err)
		return
	}
}

func (server *Server) FilterClientByLastNameAndFirstName(u *[]models.User, w http.ResponseWriter, firstName, lastName string) {
	var output = []models.User{}
	for _, c := range *u {
		if strings.ToLower(c.FirstName) == strings.ToLower(firstName) {
			output = append(output, c)
		}
	}
	var finalList = []models.User{}
	if len(output) > 0 {
		for _, o := range output {
			if strings.ToLower(o.LastName) == strings.ToLower(lastName) {
				finalList = append(finalList, o)
			}
		}
		if len(finalList) > 0 {

			clientsWithoutPassword := []UserWithoutPassword{}

			for _, client := range finalList {
				newUserWithoutPassword := new(UserWithoutPassword)
				newUserWithoutPassword.ID = int(client.ID)
				newUserWithoutPassword.FirstName = client.FirstName
				newUserWithoutPassword.LastName = client.LastName
				newUserWithoutPassword.Email = client.Email
				newUserWithoutPassword.PhoneNumber = client.PhoneNumber
				newUserWithoutPassword.Role = client.Role
				clientsWithoutPassword = append(clientsWithoutPassword, *newUserWithoutPassword)
			}

			if len(finalList) == 1 {
				json.NewEncoder(w).Encode(clientsWithoutPassword[0])
			} else {
				json.NewEncoder(w).Encode(clientsWithoutPassword)
			}
		} else {
			err := errors.New("404 Error: no client found")
			exceptions.ERROR(w, http.StatusNotFound, err)
			return
		}
	}
}

func (server *Server) FilterClientByIDAndLastName(u models.User, id int, w http.ResponseWriter, lastName string) {
	client, err := u.GetClientByID(server.Database, id)
	if err != nil {
		err = errors.New("404 Error: client not found")
		exceptions.ERROR(w, http.StatusNotFound, err)
		return
	}

	if strings.ToLower(client.LastName) == strings.ToLower(lastName) {

		newUserWithoutPassword := new(UserWithoutPassword)
		newUserWithoutPassword.ID = int(client.ID)
		newUserWithoutPassword.FirstName = client.FirstName
		newUserWithoutPassword.LastName = client.LastName
		newUserWithoutPassword.Email = client.Email
		newUserWithoutPassword.PhoneNumber = client.PhoneNumber
		newUserWithoutPassword.Role = client.Role

		json.NewEncoder(w).Encode(newUserWithoutPassword)
	} else {
		err = errors.New("404 Error: client not found")
		exceptions.ERROR(w, http.StatusNotFound, err)
		return
	}
}

func (server *Server) FilterClientByIDAndFirstName(u models.User, id int, w http.ResponseWriter, firstName string) {
	client, err := u.GetClientByID(server.Database, id)
	if err != nil {
		err = errors.New("404 Error: client not found")
		exceptions.ERROR(w, http.StatusNotFound, err)
		return
	}

	if strings.ToLower(client.FirstName) == strings.ToLower(firstName) {

		newUserWithoutPassword := new(UserWithoutPassword)
		newUserWithoutPassword.ID = int(client.ID)
		newUserWithoutPassword.FirstName = client.FirstName
		newUserWithoutPassword.LastName = client.LastName
		newUserWithoutPassword.Email = client.Email
		newUserWithoutPassword.PhoneNumber = client.PhoneNumber
		newUserWithoutPassword.Role = client.Role

		json.NewEncoder(w).Encode(newUserWithoutPassword)
	} else {
		err = errors.New("404 Error: client not found")
		exceptions.ERROR(w, http.StatusNotFound, err)
		return
	}
}

func (server *Server) FilterClientByIDFirstNameAndLastName(u models.User, id int, w http.ResponseWriter, firstName, lastName string) {
	client, err := u.GetClientByID(server.Database, id)
	if err != nil {
		err = errors.New("404 Error: client not found")
		exceptions.ERROR(w, http.StatusNotFound, err)
		return
	}

	if strings.ToLower(client.FirstName) == strings.ToLower(firstName) && strings.ToLower(client.LastName) == strings.ToLower(lastName) {

		newUserWithoutPassword := new(UserWithoutPassword)
		newUserWithoutPassword.ID = int(client.ID)
		newUserWithoutPassword.FirstName = client.FirstName
		newUserWithoutPassword.LastName = client.LastName
		newUserWithoutPassword.Email = client.Email
		newUserWithoutPassword.PhoneNumber = client.PhoneNumber
		newUserWithoutPassword.Role = client.Role

		json.NewEncoder(w).Encode(client)
	} else {
		err = errors.New("404 Error: client not found")
		exceptions.ERROR(w, http.StatusNotFound, err)
		return
	}
}

func (server *Server) FilterClientByIDWithLastNameContains(u models.User, id int, w http.ResponseWriter, c string) {
	client, err := u.GetClientByID(server.Database, id)
	if err != nil {
		err = errors.New("404 Error: client not found")
		exceptions.ERROR(w, http.StatusNotFound, err)
		return
	}
	if strings.Contains(strings.ToLower(client.LastName), strings.ToLower(c)) {

		newUserWithoutPassword := new(UserWithoutPassword)
		newUserWithoutPassword.ID = int(client.ID)
		newUserWithoutPassword.FirstName = client.FirstName
		newUserWithoutPassword.LastName = client.LastName
		newUserWithoutPassword.Email = client.Email
		newUserWithoutPassword.PhoneNumber = client.PhoneNumber
		newUserWithoutPassword.Role = client.Role

		json.NewEncoder(w).Encode(client)
	} else {
		err = errors.New("404 Error: client not found")
		exceptions.ERROR(w, http.StatusNotFound, err)
		return
	}
}

func (server *Server) FilterClientByIDWithFirstNameContains(u models.User, id int, w http.ResponseWriter, c string) {
	client, err := u.GetClientByID(server.Database, id)
	if err != nil {
		err = errors.New("404 Error: client not found")
		exceptions.ERROR(w, http.StatusNotFound, err)
		return
	}
	if strings.Contains(strings.ToLower(client.FirstName), strings.ToLower(c)) {

		newUserWithoutPassword := new(UserWithoutPassword)
		newUserWithoutPassword.ID = int(client.ID)
		newUserWithoutPassword.FirstName = client.FirstName
		newUserWithoutPassword.LastName = client.LastName
		newUserWithoutPassword.Email = client.Email
		newUserWithoutPassword.PhoneNumber = client.PhoneNumber
		newUserWithoutPassword.Role = client.Role

		json.NewEncoder(w).Encode(client)
	} else {
		err = errors.New("404 Error: client not found")
		exceptions.ERROR(w, http.StatusNotFound, err)
		return
	}
}

func (server *Server) FilterClientByIDWithFirstNameContainsAndLastNameContains(u models.User, id int, w http.ResponseWriter, lastName, firstName string) {
	client, err := u.GetClientByID(server.Database, id)
	if err != nil {
		err = errors.New("404 Error: client not found")
		exceptions.ERROR(w, http.StatusNotFound, err)
		return
	}
	if strings.Contains(strings.ToLower(client.FirstName), strings.ToLower(firstName)) && strings.Contains(strings.ToLower(client.LastName), strings.ToLower(lastName)) {

		newUserWithoutPassword := new(UserWithoutPassword)
		newUserWithoutPassword.ID = int(client.ID)
		newUserWithoutPassword.FirstName = client.FirstName
		newUserWithoutPassword.LastName = client.LastName
		newUserWithoutPassword.Email = client.Email
		newUserWithoutPassword.PhoneNumber = client.PhoneNumber
		newUserWithoutPassword.Role = client.Role

		json.NewEncoder(w).Encode(client)
	} else {
		err = errors.New("404 Error: client not found")
		exceptions.ERROR(w, http.StatusNotFound, err)
		return
	}
}

func (server *Server) FilterClientWithLastNameContains(u *[]models.User, w http.ResponseWriter, lastName string) {
	var output = []models.User{}
	for _, c := range *u {
		if strings.Contains(strings.ToLower(c.LastName), strings.ToLower(lastName)) {
			output = append(output, c)
		}
	}
	if len(output) > 0 {

		clientsWithoutPassword := []UserWithoutPassword{}

		for _, client := range output {
			newUserWithoutPassword := new(UserWithoutPassword)
			newUserWithoutPassword.ID = int(client.ID)
			newUserWithoutPassword.FirstName = client.FirstName
			newUserWithoutPassword.LastName = client.LastName
			newUserWithoutPassword.Email = client.Email
			newUserWithoutPassword.PhoneNumber = client.PhoneNumber
			newUserWithoutPassword.Role = client.Role
			clientsWithoutPassword = append(clientsWithoutPassword, *newUserWithoutPassword)
		}

		if len(output) == 1 {
			json.NewEncoder(w).Encode(clientsWithoutPassword[0])
		} else {
			json.NewEncoder(w).Encode(clientsWithoutPassword)
		}
	} else {
		err := errors.New("404 Error: no client found")
		exceptions.ERROR(w, http.StatusNotFound, err)
		return
	}
}

func (server *Server) FilterClientWithFirstNameContains(u *[]models.User, w http.ResponseWriter, firstName string) {
	var output = []models.User{}
	for _, c := range *u {
		if strings.Contains(strings.ToLower(c.FirstName), strings.ToLower(firstName)) {
			output = append(output, c)
		}
	}
	if len(output) > 0 {

		clientsWithoutPassword := []UserWithoutPassword{}

		for _, client := range output {
			newUserWithoutPassword := new(UserWithoutPassword)
			newUserWithoutPassword.ID = int(client.ID)
			newUserWithoutPassword.FirstName = client.FirstName
			newUserWithoutPassword.LastName = client.LastName
			newUserWithoutPassword.Email = client.Email
			newUserWithoutPassword.PhoneNumber = client.PhoneNumber
			newUserWithoutPassword.Role = client.Role
			clientsWithoutPassword = append(clientsWithoutPassword, *newUserWithoutPassword)
		}

		if len(output) == 1 {
			json.NewEncoder(w).Encode(clientsWithoutPassword[0])
		} else {
			json.NewEncoder(w).Encode(clientsWithoutPassword)
		}
	} else {
		err := errors.New("404 Error: no client found")
		exceptions.ERROR(w, http.StatusNotFound, err)
		return
	}
}

func (server *Server) FilterClientWithLastNameContainsAndFirstNameContains(u *[]models.User, w http.ResponseWriter, lastName, firstName string) {
	var output = []models.User{}
	for _, c := range *u {
		if strings.Contains(strings.ToLower(c.FirstName), strings.ToLower(firstName)) {
			output = append(output, c)
		}
	}
	var finalList = []models.User{}
	if len(output) > 0 {
		for _, o := range output {
			if strings.Contains(strings.ToLower(o.LastName), strings.ToLower(lastName)) {
				finalList = append(finalList, o)
			}
		}
		if len(finalList) > 0 {

			clientsWithoutPassword := []UserWithoutPassword{}

			for _, client := range finalList {
				newUserWithoutPassword := new(UserWithoutPassword)
				newUserWithoutPassword.ID = int(client.ID)
				newUserWithoutPassword.FirstName = client.FirstName
				newUserWithoutPassword.LastName = client.LastName
				newUserWithoutPassword.Email = client.Email
				newUserWithoutPassword.PhoneNumber = client.PhoneNumber
				newUserWithoutPassword.Role = client.Role
				clientsWithoutPassword = append(clientsWithoutPassword, *newUserWithoutPassword)
			}

			if len(finalList) == 1 {
				json.NewEncoder(w).Encode(clientsWithoutPassword[0])
			} else {
				json.NewEncoder(w).Encode(clientsWithoutPassword)
			}
		} else {
			err := errors.New("404 Error: no client found")
			exceptions.ERROR(w, http.StatusNotFound, err)
			return
		}
	}
}
