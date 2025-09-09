package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Sospeter-M-Kinyanjui/task-manager/database"
	"github.com/Sospeter-M-Kinyanjui/task-manager/models"
	"github.com/gorilla/mux"
)

func GetTasks(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	rows, _ := database.DB.Query(context.Background(),
		"SELECT id, title, completed FROM tasks WHERE user_id=$1", userID)

	tasks := []models.Task{}

	for rows.Next() {
		var task models.Task
		rows.Scan(&task.ID, &task.Title, &task.Completed)
		task.UserId = userID
		tasks = append(tasks, task)
	}
	json.NewEncoder(w).Encode(tasks)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	var task models.Task
	json.NewDecoder(r.Body).Decode(&task)

	err := database.DB.QueryRow(context.Background(),
		"INSERT INTO tasks (title, completed, user_id) VALUES ($1, $2, $3) RETURNING id",
		task.Title, task.Completed, userID).Scan(&task.ID)

	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	task.ID = userID
	json.NewEncoder(w).Encode(task)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var task models.Task
	json.NewDecoder(r.Body).Decode(&task)

	_, err := database.DB.Exec(context.Background(),
		"UPDATE tasks SET title=$1, completed=$2 WHERE id=$3 AND user_ID=$4",
		task.Title, task.Completed, id, userID)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	_, err := database.DB.Exec(context.Background(),
		"DELETE FROM task WHERE id=$1 AND user_id=$2", id, userID)

	if err != nil {
		http.Error(w, "DB error", http.StatusNoContent)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
