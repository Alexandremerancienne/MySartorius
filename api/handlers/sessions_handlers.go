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
	"time"

	"github.com/Alexandremerancienne/my_Sartorius/api/auth"
	"github.com/Alexandremerancienne/my_Sartorius/api/exceptions"
	"github.com/Alexandremerancienne/my_Sartorius/api/models"
	"github.com/gorilla/mux"
)

func (server *Server) CreateSession(w http.ResponseWriter, r *http.Request) {

	user := models.User{}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		err = errors.New("Invalid token: Please login with valid credentials.")
		exceptions.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	k, err := user.GetUserByID(server.Database, tokenID)
	if err != nil {
		err = errors.New("500 Internal Server Error: Cannot get user")
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
		session := models.Session{}
		if k.Role == "coach" {
			session.CoachID = int(k.ID)
		}
		if err = json.Unmarshal(body, &session); err != nil {
			err = errors.New("422 Unprocessable Entity Error: Cannot parse JSON data")
			exceptions.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}

		s := models.Session{}
		allSessions, err := s.GetSessions(server.Database)
		if err != nil {
			err = errors.New("500 Internal Server Error: Cannot get sessions")
			exceptions.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		sessionMonth := ""
		if len(strconv.Itoa(session.Month)) == 1 {
			sessionMonth = "0" + strconv.Itoa(session.Month)
		} else {
			sessionMonth = strconv.Itoa(session.Month)
		}

		sessionDay := ""
		if len(strconv.Itoa(session.Day)) == 1 {
			sessionDay = "0" + strconv.Itoa(session.Day)
		} else {
			sessionDay = strconv.Itoa(session.Day)
		}

		t3, err := time.Parse(
			time.RFC3339,
			strconv.Itoa(session.Year)+"-"+sessionMonth+
				"-"+sessionDay+"T"+session.StartingTime+":00+00:00")
		if err != nil {
			err = errors.New("400 Bad Request Error: Cannot create session with this date")
			exceptions.ERROR(w, http.StatusBadRequest, err)
			return
		}

		if t3.Format(time.RFC3339) < time.Now().Format(time.RFC3339) {
			err = errors.New("400 Bad Request Error: Cannot create session with a past date")
			exceptions.ERROR(w, http.StatusBadRequest, err)
			return
		}

		session.DateSession = strings.Replace(t3.Format(time.RFC3339)[:16], "T", " ", -1)

		t1, _ := time.Parse(
			time.RFC3339,
			strconv.Itoa(session.Year)+"-"+sessionMonth+"-"+sessionDay+"T"+session.StartingTime+":00+00:00")

		for _, sess := range *allSessions {

			sessionMonth := ""
			if len(strconv.Itoa(sess.Month)) == 1 {
				sessionMonth = "0" + strconv.Itoa(sess.Month)
			} else {
				sessionMonth = strconv.Itoa(sess.Month)
			}

			sessionDay := ""
			if len(strconv.Itoa(sess.Day)) == 1 {
				sessionDay = "0" + strconv.Itoa(sess.Day)
			} else {
				sessionDay = strconv.Itoa(sess.Day)
			}

			t2, _ := time.Parse(
				time.RFC3339,
				strconv.Itoa(sess.Year)+"-"+sessionMonth+"-"+sessionDay+"T"+sess.StartingTime+":00+00:00")

			if (t1.Format(time.RFC3339) >= t2.Format(time.RFC3339) &&
				t1.Format(time.RFC3339) <= t2.Add(time.Duration(sess.Duration)*time.Minute).Format(time.RFC3339)) &&
				(sess.ClientID == session.ClientID || sess.CoachID == session.CoachID) {
				err = errors.New("400 Bad Request Error: This slot is not available. Please choose another date")
				exceptions.ERROR(w, http.StatusBadRequest, err)
				return
			}
		}

		newSession, err := session.CreateSession(server.Database)
		if err != nil {
			err = errors.New("400 Bad Request Error: A session cannot be created for these users. Please check the IDs provided")
			exceptions.ERROR(w, http.StatusBadRequest, err)
			return
		}

		server.Database.Save(newSession)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(newSession)

	} else {
		err = errors.New("401 Error: Access restricted to Management and coaches")
		exceptions.ERROR(w, http.StatusUnauthorized, err)
		return
	}

}

func (server *Server) GetSessions(w http.ResponseWriter, r *http.Request) {

	u, err := url.Parse(r.URL.String())
	if err != nil {
		err = errors.New("422 Unprocessable Entity Error: Cannot parse url")
		exceptions.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	qs := u.Query()
	y := server.readInt(qs, "coach_id", 0)
	z := server.readInt(qs, "client_id", 0)

	sessions := models.Session{}
	retrievedSessions, err := sessions.GetSessions(server.Database)
	if err != nil {
		err = errors.New("500 Internal Server Error: Cannot get sessions")
		exceptions.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	sessionsSlice := []models.Session{}

	if len(*retrievedSessions) == 0 {

		user := models.User{}
		tokenID, err := auth.ExtractTokenID(r)
		if err != nil {
			err = errors.New("Invalid token: Please login with valid credentials.")
			exceptions.ERROR(w, http.StatusUnauthorized, err)
			return
		}

		k, err := user.GetUserByID(server.Database, tokenID)
		if err != nil {
			err = errors.New("500 Internal Server Error: Cannot get user")
			exceptions.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		if tokenID == uint32(k.ID) &&
			(k.Role == "manager" || (k.Role == "coach") || (k.Role == "client")) {
		} else {
			err = errors.New("Missing credentials: Access restricted to Management, session coach and session client.")
			exceptions.ERROR(w, http.StatusUnauthorized, err)
			return
		}

	} else {
		for _, s := range *retrievedSessions {

			u := models.User{}
			retrievedCoach, err := u.GetCoachByID(server.Database, s.CoachID)

			h := models.User{}
			retrievedClient, err := h.GetClientByID(server.Database, s.ClientID)

			user := models.User{}
			tokenID, err := auth.ExtractTokenID(r)
			if err != nil {
				err = errors.New("Invalid token: Please login with valid credentials.")
				exceptions.ERROR(w, http.StatusUnauthorized, err)
				return
			}

			k, err := user.GetUserByID(server.Database, tokenID)
			if err != nil {
				err = errors.New("500 Internal Server Error: Cannot get user")
				exceptions.ERROR(w, http.StatusInternalServerError, err)
				return
			}

			if tokenID == uint32(k.ID) &&
				(k.Role == "manager" ||
					(k.Role == "coach" && k.ID == retrievedCoach.ID) ||
					(k.Role == "client" && k.ID == retrievedClient.ID)) {

				sessionsSlice = append(sessionsSlice, s)
			} else if k.Role == "coach" {
			} else if tokenID == uint32(k.ID) && k.Role == "" {
				err = errors.New("Missing credentials: Access restricted to Management, session coach and session client.")
				exceptions.ERROR(w, http.StatusUnauthorized, err)
				return
			}
		}

		s := models.Session{}

		if y != 0 || z != 0 {
			w.Header().Set("Content-Type", "application/json")
			sessions, err := s.GetSessions(server.Database)
			if err != nil {
				err = errors.New("500 Internal Server Error: cannot get sessions")
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

			server.ApplySessionFilters(w, sessions, y, z, *j)
		} else {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(sessionsSlice)
		}
	}

}

func (server *Server) GetSession(w http.ResponseWriter, r *http.Request) {

	session := models.Session{}
	params := mux.Vars(r)
	uid, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		err = errors.New("422 Unprocessable Entity Error: Cannot parse url")
		exceptions.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	retrievedSession, err := session.GetSessionByID(server.Database, uint32(uid))

	user := models.User{}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		err = errors.New("Invalid token: Please login with valid credentials.")
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
			(k.Role == "coach" && int(k.ID) == retrievedSession.CoachID) ||
			(k.Role == "client" && k.ID == uint(retrievedSession.ClientID))) {
		retrievedSession, err := session.GetSessionByID(server.Database, uint32(uid))
		if err != nil {
			err = errors.New("500 Internal Server Error: Cannot get session")
			exceptions.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(retrievedSession)
	} else {
		err = errors.New("Missing credentials: Access restricted to Management, session coach and session client.")
		exceptions.ERROR(w, http.StatusUnauthorized, err)
		return
	}
}

func (server *Server) UpdateSession(w http.ResponseWriter, r *http.Request) {

	session := models.Session{}
	params := mux.Vars(r)
	uid, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		err = errors.New("422 Unprocessable Entity Error: Cannot parse url")
		exceptions.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	retrievedSession, err := session.GetSessionByID(server.Database, uint32(uid))

	user := models.User{}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		err = errors.New("Invalid token: Please login with valid credentials.")
		exceptions.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	k, err := user.GetUserByID(server.Database, tokenID)
	if err != nil {
		err = errors.New("500 Internal Server Error: Cannot get user")
		exceptions.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if tokenID == uint32(k.ID) && (k.Role == "manager" || (k.Role == "coach" && k.ID == uint(retrievedSession.CoachID))) {

		sessionToUpdate, err := session.GetSessionByID(server.Database, uint32(uid))
		if err != nil {
			err = errors.New("500 Internal Server Error: Cannot get session")
			exceptions.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			exceptions.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}

		err = json.Unmarshal(body, &sessionToUpdate)
		if err != nil {
			err = errors.New("422 Unprocessable Entity Error: Cannot parse JSON data")
			exceptions.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}

		sessionToUpdateMonth := ""
		if len(strconv.Itoa(sessionToUpdate.Month)) == 1 {
			sessionToUpdateMonth = "0" + strconv.Itoa(sessionToUpdate.Month)
		} else {
			sessionToUpdateMonth = strconv.Itoa(sessionToUpdate.Month)
		}

		sessionToUpdateDay := ""
		if len(strconv.Itoa(sessionToUpdate.Day)) == 1 {
			sessionToUpdateDay = "0" + strconv.Itoa(sessionToUpdate.Day)
		} else {
			sessionToUpdateDay = strconv.Itoa(sessionToUpdate.Day)
		}

		t1, err := time.Parse(
			time.RFC3339,
			strconv.Itoa(sessionToUpdate.Year)+"-"+sessionToUpdateMonth+"-"+sessionToUpdateDay+
				"T"+sessionToUpdate.StartingTime+":00+00:00")

		if err != nil {
			err = errors.New("400 Bad Request Error: Cannot update session with this date")
			exceptions.ERROR(w, http.StatusBadRequest, err)
			return
		}

		sessionToUpdate.DateSession = strings.Replace(t1.Format(time.RFC3339)[:16], "T", " ", -1)

		s := models.Session{}

		allSessions, err := s.GetSessions(server.Database)
		if err != nil {
			err = errors.New("500 Internal Server Error: Cannot get sessions")
			exceptions.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		for _, singleSession := range *allSessions {
			if singleSession.ID == sessionToUpdate.ID {

				for _, sess := range *allSessions {

					if sess != singleSession {

						sessionMonth := ""
						if len(strconv.Itoa(sess.Month)) == 1 {
							sessionMonth = "0" + strconv.Itoa(sess.Month)
						} else {
							sessionMonth = strconv.Itoa(sess.Month)
						}

						sessionDay := ""
						if len(strconv.Itoa(sess.Day)) == 1 {
							sessionDay = "0" + strconv.Itoa(sess.Day)
						} else {
							sessionDay = strconv.Itoa(sess.Day)
						}

						t2, _ := time.Parse(
							time.RFC3339,
							strconv.Itoa(sess.Year)+"-"+sessionMonth+"-"+sessionDay+"T"+sess.StartingTime+":00+00:00")

						if (t1.Format(time.RFC3339) >= t2.Format(time.RFC3339) &&
							t1.Format(time.RFC3339) <= t2.Add(time.Duration(sess.Duration)*time.Minute).Format(time.RFC3339)) &&
							(sess.ClientID == session.ClientID || sess.CoachID == session.CoachID) {
							err = errors.New("400 Bad Request Error: This slot is not available. Please choose another date")
							exceptions.ERROR(w, http.StatusBadRequest, err)
							return
						}
					}
				}

				singleSession.Title = sessionToUpdate.Title
				singleSession.Description = sessionToUpdate.Description
				singleSession.Year = sessionToUpdate.Year
				singleSession.Month = sessionToUpdate.Month
				singleSession.Day = sessionToUpdate.Day


				singleSession.StartingTime = sessionToUpdate.StartingTime
				singleSession.Duration = sessionToUpdate.Duration
				singleSession.DateSession = sessionToUpdate.DateSession


				server.Database.Save(singleSession)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(singleSession)
			}
		}

	} else {
		err = errors.New("Missing credentials: Access restricted to Management and session coach.")
		exceptions.ERROR(w, http.StatusUnauthorized, err)
		return

	}

	w.Header().Set("Content-Type", "application/json")
}

func (server *Server) DeleteSession(w http.ResponseWriter, r *http.Request) {

	session := models.Session{}
	params := mux.Vars(r)
	uid, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		err = errors.New("422 Unprocessable Entity Error: Cannot parse url")
		exceptions.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	retrievedSession, err := session.GetSessionByID(server.Database, uint32(uid))

	user := models.User{}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		err = errors.New("Invalid token: Please login with valid credentials.")
		exceptions.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	k, err := user.GetUserByID(server.Database, tokenID)
	if err != nil {
		err = errors.New("500 Internal Server Error: Cannot get user")
		exceptions.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if tokenID == uint32(k.ID) && (k.Role == "manager" || (k.Role == "coach" && k.ID == uint(retrievedSession.CoachID))) {
		if _, err = session.DeleteSession(server.Database, uint32(uid)); err != nil {
			err = errors.New("500 Internal Server Error: Cannot delete session")
			exceptions.ERROR(w, http.StatusInternalServerError, err)
			return
		}
	} else {
		err = errors.New("Missing credentials: Access restricted to Management and session coach.")
		exceptions.ERROR(w, http.StatusUnauthorized, err)
		return

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(204)
	er := json.NewEncoder(w).Encode("")
	if er != nil {
		fmt.Fprintf(w, "%s", er.Error())
	}
}

