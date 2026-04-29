package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"todo-api/model"
	"todo-api/queue"

	"github.com/gorilla/mux"
)

type TodoHandler struct {
	DB *sql.DB
	MQ *queue.RabbitMQ
}

func NewTodoHandler(db *sql.DB, mq *queue.RabbitMQ) *TodoHandler {
	return &TodoHandler{DB: db, MQ: mq}
}

func (h *TodoHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query("SELECT id, title, completed FROM todos")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	todos := []model.Todo{}
	for rows.Next() {
		var t model.Todo
		rows.Scan(&t.ID, &t.Title, &t.Completed)
		todos = append(todos, t)
	}

	w.Header().Set("Content-Type", "applicationjson")
	json.NewEncoder(w).Encode(todos)
}

func (h *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var t model.Todo
	json.NewDecoder(r.Body).Decode(&t)

	err := h.DB.QueryRow("INSERT INTO todos (title, completed) VALUES ($1, $2) RETURNING id", t.Title, t.Completed).Scan(&t.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if h.MQ != nil {
		msg := fmt.Sprintf("Created a new task: %s", t.Title)
		h.MQ.Publish("notifications", msg)
	}

	w.Header().Set("Content-Type", "applicationjson")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}

func (h *TodoHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var t model.Todo
	json.NewDecoder(r.Body).Decode(&t)

	_, err := h.DB.Exec("UPDATE todos SET title=$1, completed=$2 WHERE id=$3", t.Title, t.Completed, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t.ID = id
	w.Header().Set("Content-Type", "applicationjson")
	json.NewEncoder(w).Encode(t)
}

func (h *TodoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	_, err := h.DB.Exec("DELETE FROM todos where id=$1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
