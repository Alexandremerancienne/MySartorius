package seed

import (
	"fmt"
	"log"

	"github.com/Alexandremerancienne/my_Sartorius/api/handlers"
	"github.com/Alexandremerancienne/my_Sartorius/api/models"
)

func RefreshTables(s handlers.Server) error {

	err := s.Database.Migrator().DropTable("users", "clients", "coaches", "sessions", "tasks", "reminders")
	if err != nil {
		fmt.Println("Cannot refresh database tables")
		log.Fatal("Error:", err)
	}

	err = s.Database.AutoMigrate(&models.User{}, &models.Coach{}, &models.Client{}, &models.Session{}, &models.Task{}, &models.Reminder{})
	if err != nil {
		fmt.Println("Cannot migrate data")
		log.Fatal("Error:", err)
	}
	log.Printf("Successfully refreshed tables")
	return nil
}

func SeedTables(s handlers.Server) ([]models.User, []models.Coach, []models.Client, []models.Session, []models.Task, []models.Reminder, error) {

	RefreshTables(s)

	var err error

	var users = []models.User{
		{
			FirstName:   "Raphaël",
			LastName:    "Merancienne",
			Email:       "raf@bot.com",
			Password:    "admin",
			PhoneNumber: "+33497854010",
			Role:        "manager",
		},
		{
			FirstName:   "Maria",
			LastName:    "Petrova",
			Email:       "maria_coach@bot.com",
			Password:    "hellogophers",
			PhoneNumber: "+33646452149",
			Role:        "coach",
		},
		{
			FirstName:   "Dimitar",
			LastName:    "Penchev",
			Email:       "dimitar_client@bot.com",
			Password:    "hellogophers",
			PhoneNumber: "+33162481012",
			Role:        "client",
		},
		{
			FirstName:   "Emil",
			LastName:    "Donchev",
			Email:       "emil_coach@bot.com",
			Password:    "hellogophers",
			PhoneNumber: "+33559865840",
			Role:        "coach",
		},
		{
			FirstName:   "Irina",
			LastName:    "Angelova",
			Email:       "irina_client@bot.com",
			Password:    "hellogophers",
			PhoneNumber: "+3355978135",
			Role:        "client",
		},
	}

	var coaches = []models.Coach{
		{
			ID:     2,
			UserID: 2,
		},
		{
			ID:     4,
			UserID: 4,
		},
	}

	var clients = []models.Client{
		{
			ID:     3,
			UserID: 3,
		},
		{
			ID:     5,
			UserID: 5,
		},
	}

	var sessions = []models.Session{
		{
			Title:        "First session Dimitar",
			Description:  "Intro session: cardio/workout/abs",
			CoachID:      2,
			ClientID:     3,
			Year:         2022,
			Month:        11,
			Day:          21,
			StartingTime: "15:00",
			Duration:     60,
			DateSession:  "2022-11-21 15:00",
		},
		{
			Title:        "Second session Dimitar",
			Description:  "Follow-up session: cardio/workout/abs",
			CoachID:      2,
			ClientID:     3,
			Year:         2022,
			Month:        12,
			Day:          15,
			StartingTime: "15:00",
			Duration:     60,
			DateSession:  "2022-12-05 15:00",
		},
		{
			Title:        "First session Irina",
			Description:  "Intro session: cardio/workout/abs",
			CoachID:      4,
			ClientID:     5,
			Year:         2022,
			Month:        11,
			Day:          21,
			StartingTime: "15:00",
			Duration:     60,
			DateSession:  "2022-11-21 15:00",
		},
		{
			Title:        "Second session Irina",
			Description:  "Follow-up session: cardio/workout/abs",
			CoachID:      4,
			ClientID:     5,
			Year:         2022,
			Month:        12,
			Day:          15,
			StartingTime: "15:00",
			Duration:     60,
			DateSession:  "2022-12-05 15:00",
		},
	}

	var tasks = []models.Task{
		{
			Title:       "Monday workout",
			Description: "10 abs, 15 push-ups, 10 burpees. 30 seconds rest between each exercise",
			AssignerID:  2,
			AssigneeID:  3,
			Year:        2022,
			Month:       11,
			Day:         15,
			Duration:    30,
			DateTask:    "2022-11-15",
		},
		{
			Title:       "Wednesday workout",
			Description: "15 abs, 20 push-ups, 15 burpees. 35 seconds rest between each exercise",
			AssignerID:  2,
			AssigneeID:  3,
			Year:        2022,
			Month:       11,
			Day:         17,
			Duration:    30,
			DateTask:    "2022-11-17",
		},
		{
			Title:       "Tuesday workout",
			Description: "5 abs, 5 push-ups, 5 burpees. 40 seconds rest between each exercise",
			AssignerID:  4,
			AssigneeID:  5,
			Year:        2022,
			Month:       11,
			Day:         16,
			Duration:    30,
			DateTask:    "2022-11-16",
		},
		{
			Title:       "Thursday workout",
			Description: "7 abs, 8 push-ups, 8 burpees. 45 seconds rest between each exercise",
			AssignerID:  4,
			AssigneeID:  5,
			Year:        2022,
			Month:       11,
			Day:         18,
			Duration:    30,
			DateTask:    "2022-11-18",
		},
	}

	var reminders = []models.Reminder{
		{
			TaskID:      1,
			Description: "Do not forget to bend the arms on the push-ups",
		},
		{
			TaskID:      4,
			Description: "Drink between the exercises and not only at the end",
		},
		{
			TaskID:      4,
			Description: "Use timer for the burpees",
		},
	}

	for i, _ := range users {
		err = s.Database.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("Cannot seed users table: %v", err)
		}
	}

	for i, _ := range coaches {
		err = s.Database.Model(&models.Coach{}).Create(&coaches[i]).Error
		if err != nil {
			log.Fatalf("Cannot seed coaches table: %v", err)
		}
	}

	for i, _ := range clients {
		err = s.Database.Model(&models.Client{}).Create(&clients[i]).Error
		if err != nil {
			log.Fatalf("Cannot seed clients table: %v", err)
		}
	}

	for i, _ := range sessions {
		err = s.Database.Model(&models.Session{}).Create(&sessions[i]).Error
		if err != nil {
			log.Fatalf("Cannot seed sessions table: %v", err)
		}
	}

	for i, _ := range tasks {
		err = s.Database.Model(&models.Task{}).Create(&tasks[i]).Error
		if err != nil {
			log.Fatalf("Cannot seed tasks table: %v", err)
		}
	}

	for i, _ := range reminders {
		err = s.Database.Model(&models.Reminder{}).Create(&reminders[i]).Error
		if err != nil {
			log.Fatalf("Cannot seed reminders table: %v", err)
		}
	}

	return users, coaches, clients, sessions, tasks, reminders, nil
}
