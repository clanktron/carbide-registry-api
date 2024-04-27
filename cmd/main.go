package main

import (
	"carbide-registry-api/pkg/api"
	"carbide-registry-api/pkg/database"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	db_user := os.Getenv("DBUSER")
	if db_user == "" {
		db_user = "clayton"
	}
	db_pass := os.Getenv("DBPASS")
	if db_pass == "" {
		db_pass = "applevisioncurescancer"
	}
	db_host := os.Getenv("DBHOST")
	if db_host == "" {
		db_host = "127.0.0.1"
	}
	db_port := os.Getenv("DBPORT")
	if db_port == "" {
		db_port = "3306"
	}
	db_name := os.Getenv("DBNAME")
	if db_name == "" {
		db_name = "carbide"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	db, err := database.Init(db_user, db_pass, db_host, db_port, db_name)
	if err != nil {
		log.Error("Database connection failed, exiting...")
		log.Fatal(err)
	}
	log.Info("Database connected!")
	defer db.Close()
	log.Info("Initializing db schema...")
	err = database.SchemaInit(db)
	if err != nil {
		log.Error("Database schema init failed, exiting...")
		log.Fatal(err)
	}
	log.Info("Starting server on port " + port + "...")
	log.Fatal(http.ListenAndServe(":"+port, api.NewRouter(db)))
}
