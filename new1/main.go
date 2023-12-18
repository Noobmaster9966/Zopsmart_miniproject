package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"gofr.dev/pkg/gofr"
	"new1/pkg/handlers"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", dbUser, dbPassword, dbName))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := gofr.New()

	handlers.RegisterUserHandlers(app, db)

	app.Start()
}
