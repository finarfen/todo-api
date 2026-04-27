package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func Connect() *sql.DB {
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=todos sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("БД не доступна:", err)
	}

	fmt.Println("Подключено к PostrgeSQL!")
	return db
}
