package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Alexandremerancienne/my_Sartorius/api/handlers"
	"github.com/Alexandremerancienne/my_Sartorius/seed"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	fmt.Println("My Sartorius")

	s := handlers.Server{Router: mux.NewRouter()}

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env")
	}
	db_password := os.Getenv("DB_PASSWORD")
	db_name := os.Getenv("DB_NAME")

	dsn := "root:" + db_password + "@tcp(127.0.0.1:3306)/" + db_name + "?charset=utf8mb4&parseTime=True&loc=Local"
	s.Database, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		if err := handlers.CreateDatabase(db_password, db_name); err != nil {
			fmt.Println("Cannot connect to database")
			log.Fatal("Error:", err)
		}
	}

	seed.SeedTables(s)

	sqlDB, err := s.Database.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	fmt.Println("Connected to database")

	log.Println("Listening to Port :4000")
	s.InitializeRoutes()
}
