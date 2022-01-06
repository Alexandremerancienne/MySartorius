package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Alexandremerancienne/my_Sartorius/api/auth"
	"github.com/Alexandremerancienne/my_Sartorius/api/exceptions"
	"github.com/Alexandremerancienne/my_Sartorius/api/models"
	"github.com/gorilla/mux"
)

type UserWithoutPassword struct {
	ID int
	FirstName   string
	LastName    string
	Email       string
	PhoneNumber string
	Role        string
}

func (server *Server) readString(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}
	return s
}

func (server *Server) readCSV(qs url.Values, key string, defaultValue []string) []string {
	csv := qs.Get(key)
	if csv == "" {
		return defaultValue
	}
	return strings.Split(csv, ",")
}

func (server *Server) readInt(qs url.Values, key string, defaultValue int) int {

	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		errors.New("must be an integer value")
		return defaultValue
	}
	return i
}

func (server *Server) CheckIfEmailIsUnique(w http.ResponseWriter, email string) bool {

	user := models.User{}

	allUsers, err := user.GetUsers(server.Database)
	if err != nil {
		err = errors.New("500 Internal Server Error: cannot get users")
		exceptions.ERROR(w, http.StatusInternalServerError, err)
		return false
	}

	for _, user := range *allUsers {
		if user.Email == email {
			return false
		}
	}
	return true
}

func (server *Server) CheckIfUserIsManager(w http.ResponseWriter, r *http.Request) bool {

	user := models.User{}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		err = errors.New("Invalid token: please login with valid credentials")
		exceptions.ERROR(w, http.StatusUnauthorized, err)
		return false
	}

	k, err := user.GetUserByID(server.Database, tokenID)
	if err != nil {
		err = errors.New("500 Internal Server Error: cannot get user")
		exceptions.ERROR(w, http.StatusInternalServerError, err)
		return false
	}

	if tokenID != uint32(k.ID) || k.Role != "manager" {
		err = errors.New("Missing credentials: access restricted to Management")
		exceptions.ERROR(w, http.StatusUnauthorized, err)
		return false
	}
	return true
}

func (server *Server) GetAllUsersIfManagerOrReturnOwnProfile(w http.ResponseWriter, r *http.Request) bool {

	user := models.User{}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		err = errors.New("Invalid token: please login with valid credentials")
		exceptions.ERROR(w, http.StatusUnauthorized, err)
		return false
	}

	k, err := user.GetUserByID(server.Database, tokenID)
	if err != nil {
		err = errors.New("500 Internal Server Error: cannot get user")
		exceptions.ERROR(w, http.StatusInternalServerError, err)
		return false
	}

	if tokenID != uint32(k.ID) || k.Role != "manager" {

		newUserWithoutPassword := new(UserWithoutPassword)
		newUserWithoutPassword.ID = int(k.ID)
		newUserWithoutPassword.FirstName = k.FirstName
		newUserWithoutPassword.LastName = k.LastName
		newUserWithoutPassword.Email = k.Email
		newUserWithoutPassword.PhoneNumber = k.PhoneNumber
		newUserWithoutPassword.Role = k.Role

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(newUserWithoutPassword)
		return false
	}
	return true
}

func (server *Server) GetUserIfManagerOrReturnOwnProfile(w http.ResponseWriter, r *http.Request) bool {

	user := models.User{}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		err = errors.New("Invalid token: please login with valid credentials")
		exceptions.ERROR(w, http.StatusUnauthorized, err)
		return false
	}

	k, err := user.GetUserByID(server.Database, tokenID)
	if err != nil {
		err = errors.New("500 Internal Server Error: cannot get user")
		exceptions.ERROR(w, http.StatusInternalServerError, err)
		return false
	}

	if tokenID != uint32(k.ID) || k.Role != "manager" {

		params := mux.Vars(r)
		uid, _ := strconv.ParseUint(params["id"], 10, 32)

		if uid == uint64(k.ID) {

			newUserWithoutPassword := new(UserWithoutPassword)
			newUserWithoutPassword.ID = int(k.ID)
			newUserWithoutPassword.FirstName = k.FirstName
			newUserWithoutPassword.LastName = k.LastName
			newUserWithoutPassword.Email = k.Email
			newUserWithoutPassword.PhoneNumber = k.PhoneNumber
			newUserWithoutPassword.Role = k.Role

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(newUserWithoutPassword)
			return false
		} else {
			err = errors.New("Missing credentials: access restricted to Management")
			exceptions.ERROR(w, http.StatusUnauthorized, err)
			return false
		}

	}
	return true
}

func (server *Server) CreateUser(w http.ResponseWriter, r *http.Request) {

	token := auth.ExtractToken(r)

	if token == "" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			exceptions.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}
		user := models.User{}
		if err = json.Unmarshal(body, &user); err != nil {
			err = errors.New("422 Unprocessable Entity Error: Cannot parse JSON data")
			exceptions.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}

		if user.LastName == "" || user.FirstName == "" || user.Email == "" {
			err = errors.New("400 Error: Last name, first name and email fields must be filled")
			exceptions.ERROR(w, http.StatusBadRequest, err)
			return

		}

		e := server.CheckIfEmailIsUnique(w, user.Email)
		if !e {
			err = errors.New("400 Error: this email address is already used. Please choose another email address")
			exceptions.ERROR(w, http.StatusBadRequest, err)
			return
		}

		newUser, err := user.CreateUser(server.Database)
		if err != nil {
			err = errors.New("500 Internal Server Error: cannot create user")
			exceptions.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		newUser.Role = ""

		newUserWithoutPassword := new(UserWithoutPassword)
		newUserWithoutPassword.ID = int(newUser.ID)
		newUserWithoutPassword.FirstName = newUser.FirstName
		newUserWithoutPassword.LastName = newUser.LastName
		newUserWithoutPassword.Email = newUser.Email
		newUserWithoutPassword.PhoneNumber = newUser.PhoneNumber
		newUserWithoutPassword.Role = newUser.Role

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(newUserWithoutPassword)
	} else {
		e := server.CheckIfUserIsManager(w, r)
		if e == true {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				exceptions.ERROR(w, http.StatusUnprocessableEntity, err)
				return
			}
			user := models.User{}
			if err = json.Unmarshal(body, &user); err != nil {
				err = errors.New("422 Unprocessable Entity Error: Cannot parse JSON data")
				exceptions.ERROR(w, http.StatusUnprocessableEntity, err)
				return
			}

			if user.LastName == "" || user.FirstName == "" || user.Email == "" {
				err = errors.New("400 Error: Last name, first name and email fields must be filled")
				exceptions.ERROR(w, http.StatusBadRequest, err)
				return

			}

			e := server.CheckIfEmailIsUnique(w, user.Email)
			if !e {
				err = errors.New("400 Error: this email address is already used. Please choose another email address")
				exceptions.ERROR(w, http.StatusBadRequest, err)
				return
			}

			newUser, err := user.CreateUser(server.Database)
			if err != nil {
				err = errors.New("500 Internal Server Error: cannot create user")
				exceptions.ERROR(w, http.StatusInternalServerError, err)
				return
			}

			if newUser.Role != "client" && newUser.Role != "coach" {
				err = errors.New("422 Unprocessable Entity Error: please choose role between 'client' and 'coach' (role must be lowercase)")
				exceptions.ERROR(w, http.StatusUnprocessableEntity, err)
				return
			}

			newUserWithoutPassword := new(UserWithoutPassword)
			newUserWithoutPassword.ID = int(newUser.ID)
			newUserWithoutPassword.FirstName = newUser.FirstName
			newUserWithoutPassword.LastName = newUser.LastName
			newUserWithoutPassword.Email = newUser.Email
			newUserWithoutPassword.PhoneNumber = newUser.PhoneNumber
			newUserWithoutPassword.Role = newUser.Role

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(201)
			json.NewEncoder(w).Encode(newUserWithoutPassword)
		}
		return
	}
}

func (server *Server) GetUsers(w http.ResponseWriter, r *http.Request) {

	u, err := url.Parse(r.URL.String())
	if err != nil {
		err = errors.New("422 Unprocessable Entity Error: Cannot parse url")
		exceptions.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	qs := u.Query()
	x := server.readString(qs, "role", "")
	y := server.readInt(qs, "id", 0)
	z1 := server.readString(qs, "last_name", "")
	z2 := server.readString(qs, "first_name", "")
	z3 := server.readString(qs, "last_name_contains", "")
	z4 := server.readString(qs, "first_name_contains", "")

	boolean := server.GetAllUsersIfManagerOrReturnOwnProfile(w, r)
	if boolean != true {
		return
	}

	user := models.User{}

	switch x {
	case "coach":
		w.Header().Set("Content-Type", "application/json")

		coaches, err := user.GetCoaches(server.Database)
		if err != nil {
			err = errors.New("500 Internal Server Error: cannot get coaches")
			exceptions.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		server.ApplyCoachFilters(w, coaches, user, y, z1, z2, z3, z4)

	case "client":
		w.Header().Set("Content-Type", "application/json")

		clients, err := user.GetClients(server.Database)
		if err != nil {
			err = errors.New("500 Internal Server Error: cannot get clients")
			exceptions.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		server.ApplyClientFilters(w, clients, user, y, z1, z2, z3, z4)

	default:
		w.Header().Set("Content-Type", "application/json")
		users, err := user.GetUsers(server.Database)
		if err != nil {
			err = errors.New("500 Internal Server Error: cannot get users")
			exceptions.ERROR(w, http.StatusInternalServerError, err)
			return
		} else {
			usersWithoutPassword := []UserWithoutPassword{}
			for _, user := range *users {
				newUserWithoutPassword := new(UserWithoutPassword)
				newUserWithoutPassword.ID = int(user.ID)
				newUserWithoutPassword.FirstName = user.FirstName
				newUserWithoutPassword.LastName = user.LastName
				newUserWithoutPassword.Email = user.Email
				newUserWithoutPassword.PhoneNumber = user.PhoneNumber
				newUserWithoutPassword.Role = user.Role
				usersWithoutPassword = append(usersWithoutPassword, *newUserWithoutPassword)
			}
			json.NewEncoder(w).Encode(usersWithoutPassword)
		}
	}
}

func (server *Server) GetUser(w http.ResponseWriter, r *http.Request) {

	user := models.User{}

	params := mux.Vars(r)
	uid, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		err = errors.New("Invalid ID: User not found.")
		exceptions.ERROR(w, http.StatusNotFound, err)
		return
	}

	boolean := server.GetUserIfManagerOrReturnOwnProfile(w, r)
	if boolean != true {
		return
	}

	retrievedUser, err := user.GetUserByID(server.Database, uint32(uid))
	if err != nil {
		err = errors.New("500 Internal Server Error: cannot get user")
		exceptions.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	newUserWithoutPassword := new(UserWithoutPassword)
	newUserWithoutPassword.ID = int(retrievedUser.ID)
	newUserWithoutPassword.FirstName = retrievedUser.FirstName
	newUserWithoutPassword.LastName = retrievedUser.LastName
	newUserWithoutPassword.Email = retrievedUser.Email
	newUserWithoutPassword.PhoneNumber = retrievedUser.PhoneNumber
	newUserWithoutPassword.Role = retrievedUser.Role

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newUserWithoutPassword)
}

func (server *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {

	user := models.User{}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		err = errors.New("Invalid token: please login with valid credentials")
		exceptions.ERROR(w, http.StatusUnauthorized, err)
	}

	k, err := user.GetUserByID(server.Database, tokenID)
	if err != nil {
		err = errors.New("500 Internal Server Error: cannot get user")
		exceptions.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if !(tokenID == uint32(k.ID) && (k.Role == "manager" || k.Role == "client" || k.Role == "")) {
		err = errors.New("Missing credentials: access restricted to Management and Profile owner")
		exceptions.ERROR(w, http.StatusUnauthorized, err)
	} else {
		if tokenID == uint32(k.ID) && k.Role == "manager" {

			params := mux.Vars(r)
			userToUpdateId, err := strconv.ParseUint(params["id"], 10, 32)
			if err != nil {
				err = errors.New("422 Unprocessable Entity Error: Cannot parse url")
				exceptions.ERROR(w, http.StatusUnprocessableEntity, err)
				return
			}

			j := models.User{}
			userToUpdate, err := j.GetUserByID(server.Database, uint32(userToUpdateId))
			if err != nil {
				err = errors.New("500 Internal Server Error: cannot get user")
				exceptions.ERROR(w, http.StatusInternalServerError, err)
				return
			}

			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				exceptions.ERROR(w, http.StatusUnprocessableEntity, err)
				return
			}

			err = json.Unmarshal(body, &userToUpdate)
			if err != nil {
				err = errors.New("422 Unprocessable Entity Error: Cannot parse JSON data")
				exceptions.ERROR(w, http.StatusUnprocessableEntity, err)
				return
			}

			f := models.User{}

			allusers, err := f.GetUsers(server.Database)

			for _, singleUser := range *allusers {
				if singleUser.ID == userToUpdate.ID {
					fmt.Println(singleUser)
					if singleUser.Role == "" {
						server.Database.Model(&singleUser).Update("role", userToUpdate.Role)

						newUserWithoutPassword := new(UserWithoutPassword)
						newUserWithoutPassword.ID = int(singleUser.ID)
						newUserWithoutPassword.FirstName = singleUser.FirstName
						newUserWithoutPassword.LastName = singleUser.LastName
						newUserWithoutPassword.Email = singleUser.Email
						newUserWithoutPassword.PhoneNumber = singleUser.PhoneNumber
						newUserWithoutPassword.Role = singleUser.Role

						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(newUserWithoutPassword)
					} else {
						text := "Role already set as " + singleUser.Role +
							". If you want to change role, please delete user and add a new one."
						err = errors.New(text)
						exceptions.ERROR(w, http.StatusUnauthorized, err)
						return
					}
				}
			}

		} else if tokenID == uint32(k.ID) && (k.Role == "client" || k.Role == "") {

			params := mux.Vars(r)
			userToUpdateId, err := strconv.ParseUint(params["id"], 10, 32)
			if err != nil {
				err = errors.New("422 Unprocessable Entity Error: Cannot parse url")
				exceptions.ERROR(w, http.StatusUnprocessableEntity, err)
				return
			}

			if userToUpdateId != uint64(tokenID) {
				err = errors.New("Error: You are not allowed to update this profile")
				exceptions.ERROR(w, http.StatusUnauthorized, err)
				return
			}

			j := models.User{}
			userToUpdate, err := j.GetUserByID(server.Database, tokenID)
			if err != nil {
				err = errors.New("500 Internal Server Error: cannot get user")
				exceptions.ERROR(w, http.StatusInternalServerError, err)
				return
			}

			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				exceptions.ERROR(w, http.StatusUnprocessableEntity, err)
				return
			}

			err = json.Unmarshal(body, &userToUpdate)
			if err != nil {
				err = errors.New("422 Unprocessable Entity Error: Cannot parse JSON data")
				exceptions.ERROR(w, http.StatusUnprocessableEntity, err)
				return
			}

			fmt.Println(userToUpdate)

			f := models.User{}

			allusers, err := f.GetUsers(server.Database)

			for _, singleUser := range *allusers {
				if singleUser.ID == userToUpdate.ID {
					server.Database.Model(&singleUser).Update("first_name", userToUpdate.FirstName)
					server.Database.Model(&singleUser).Update("last_name", userToUpdate.LastName)
					server.Database.Model(&singleUser).Update("email", userToUpdate.Email)
					server.Database.Model(&singleUser).Update("phone_number", userToUpdate.PhoneNumber)
					server.Database.Model(&singleUser).Update("password", userToUpdate.Password)
					server.Database.Save(&singleUser)

					newUserWithoutPassword := new(UserWithoutPassword)
					newUserWithoutPassword.ID = int(singleUser.ID)
					newUserWithoutPassword.FirstName = singleUser.FirstName
					newUserWithoutPassword.LastName = singleUser.LastName
					newUserWithoutPassword.Email = singleUser.Email
					newUserWithoutPassword.PhoneNumber = singleUser.PhoneNumber
					newUserWithoutPassword.Role = singleUser.Role

					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(&newUserWithoutPassword)
				}
			}
		}
	}
}

func (server *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {

	err := server.CheckIfUserIsManager(w, r)
	if err != true {
		return
	}

	user := models.User{}
	params := mux.Vars(r)
	uid, p := strconv.ParseUint(params["id"], 10, 32)

	u := models.User{}
	userToDelete, e := u.GetUserByID(server.Database, uint32(uid))
	if e != nil {
		e = errors.New("500 Internal Server Error: cannot get user")
		exceptions.ERROR(w, http.StatusInternalServerError, e)
		return
	}
	if userToDelete.Role == "manager" {
		err := errors.New("401 Error: Cannot delete manager")
		exceptions.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	if userToDelete.Role == "coach" {
		server.Database.Where("id LIKE ?", uint32(uid)).Delete(models.Coach{})
	}

	if userToDelete.Role == "client" {
		server.Database.Where("id LIKE ?", uint32(uid)).Delete(models.Client{})
	}

	if p != nil {
		err := errors.New("422 Unprocessable Entity Error: Cannot parse url")
		exceptions.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	if _, c := user.DeleteUser(server.Database, uint32(uid)); c != nil {
		err := errors.New("500 Internal Server Error: cannot delete user")
		exceptions.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	s := models.Session{}
	allsessions, e := s.GetSessions(server.Database)
	if e != nil {
		e = errors.New("500 Internal Server Error: cannot get sessions")
		exceptions.ERROR(w, http.StatusInternalServerError, e)
		return
	}
	for _, session := range *allsessions {
		if session.ClientID == int(userToDelete.ID) || session.CoachID == int(userToDelete.ID) {
			i := models.Session{}
			i.DeleteSession(server.Database, uint32(session.ID))
		}
	}

	t := models.Task{}
	alltasks, o := t.GetTasks(server.Database)
	if o != nil {
		o = errors.New("500 Internal Server Error: cannot get tasks")
		exceptions.ERROR(w, http.StatusInternalServerError, e)
		return
	}
	for _, task := range *alltasks {
		if task.AssignerID == int(userToDelete.ID) || task.AssigneeID == int(userToDelete.ID) {
			i := models.Task{}
			i.DeleteTask(server.Database, uint32(task.ID))
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(204)
	er := json.NewEncoder(w).Encode("")
	if er != nil {
		fmt.Fprintf(w, "%s", er.Error())
	}
}
