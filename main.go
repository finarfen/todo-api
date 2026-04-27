package main

import (
	"log"
	"net/http"
	"todo-api/db"
	"todo-api/handler"

	"github.com/gorilla/mux"
)

func main() {
	database := db.Connect()
	defer database.Close()

	_, err := database.Exec(`
		CREATE TABLE IF NOT EXISTS todos (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			completed BOOLEAN NOT NULL DEFAULT FALSE
		)
	`)
	if err != nil {
		log.Fatal("Ошибка создания таблицы:", err)
	}

	h := handler.NewTodoHandler(database)

	r := mux.NewRouter()
	r.HandleFunc("/todos", h.GetAll).Methods("GET")
	r.HandleFunc("/todos", h.Create).Methods("POST")
	r.HandleFunc("/todos/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/todos/{id}", h.Delete).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", r))
}
