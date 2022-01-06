package handlers

import (
	"log"
	"net/http"

	"github.com/Alexandremerancienne/my_Sartorius/api/middlewares"
)

func (s *Server) InitializeRoutes() {

	s.Router.HandleFunc("/login", s.Login).Methods("POST")

	// Users routes
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareAuthentication(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareAuthentication(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/users", s.CreateUser).Methods("POST")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareAuthentication(s.UpdateUser)).Methods("PATCH")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

	// Sessions routes
	s.Router.HandleFunc("/sessions", middlewares.SetMiddlewareAuthentication(s.GetSessions)).Methods("GET")
	s.Router.HandleFunc("/sessions/{id}", middlewares.SetMiddlewareAuthentication(s.GetSession)).Methods("GET")
	s.Router.HandleFunc("/sessions", middlewares.SetMiddlewareAuthentication(s.CreateSession)).Methods("POST")
	s.Router.HandleFunc("/sessions/{id}", middlewares.SetMiddlewareAuthentication(s.UpdateSession)).Methods("PATCH")
	s.Router.HandleFunc("/sessions/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteSession)).Methods("DELETE")

	// Tasks routes
	s.Router.HandleFunc("/tasks", middlewares.SetMiddlewareAuthentication(s.GetTasks)).Methods("GET")
	s.Router.HandleFunc("/tasks/{id}", middlewares.SetMiddlewareAuthentication(s.GetTask)).Methods("GET")
	s.Router.HandleFunc("/tasks", middlewares.SetMiddlewareAuthentication(s.CreateTask)).Methods("POST")
	s.Router.HandleFunc("/tasks/{id}", middlewares.SetMiddlewareAuthentication(s.UpdateTask)).Methods("PATCH")
	s.Router.HandleFunc("/tasks/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteTask)).Methods("DELETE")

	// Reminders routes
	s.Router.HandleFunc("/tasks/{task_id}/reminders", middlewares.SetMiddlewareAuthentication(s.GetReminders)).Methods("GET")
	s.Router.HandleFunc("/tasks/{task_id}/reminders/{reminder_id}", middlewares.SetMiddlewareAuthentication(s.GetReminder)).Methods("GET")
	s.Router.HandleFunc("/tasks/{task_id}/reminders", middlewares.SetMiddlewareAuthentication(s.CreateReminder)).Methods("POST")
	s.Router.HandleFunc("/tasks/{task_id}/reminders/{reminder_id}", middlewares.SetMiddlewareAuthentication(s.UpdateReminder)).Methods("PATCH")
	s.Router.HandleFunc("/tasks/{task_id}/reminders/{reminder_id}", middlewares.SetMiddlewareAuthentication(s.DeleteReminder)).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":4000", s.Router))
}
