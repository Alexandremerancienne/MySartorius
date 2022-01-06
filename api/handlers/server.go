package handlers

import (
	"database/sql"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Server struct {
	Database *gorm.DB
	Router   *mux.Router
}

func CreateDatabase(db_password, db_name string) error {
	db, err := sql.Open("mysql", "root:"+db_password+"@tcp(127.0.0.1:3306)/")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE " + db_name)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("USE " + db_name)
	if err != nil {
		panic(err)
	}
	return nil
}
