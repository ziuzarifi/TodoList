package models

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() {

	connStr := "postgres://postgres:postgres@localhost:5432/todo_list?sslmode=disable"

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to the database")

	createTables()
}

func createTables() {
	// Create 'tasks' table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			title VARCHAR(30) UNIQUE,
			completed BOOLEAN,
			deadline TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Create 'users' table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			user_id SERIAL PRIMARY KEY,
			name VARCHAR(20) NOT NULL,
			phone_number INTEGER UNIQUE,
			email_address VARCHAR(30) UNIQUE,
			password VARCHAR(100) NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
}
