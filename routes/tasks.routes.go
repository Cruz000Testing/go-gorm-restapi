package routes

import (
	"encoding/json"
	"net/http"

	"github.com/Cruz000Testing/go-gorm-restapi/db"
	"github.com/Cruz000Testing/go-gorm-restapi/models"
	"github.com/gorilla/mux"
)

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	var tasks []models.Task
	db.DB.Find(&tasks)
	json.NewEncoder(w).Encode(tasks)
}

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	params := mux.Vars(r)
	db.DB.First(&task, params["id"])

	if task.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Task not found :O"))
		return
	}

	json.NewEncoder(w).Encode(&task)
}

func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	jsonErr := json.NewDecoder(r.Body).Decode(&task)

	if jsonErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid JSON"))
		return
	}

	createdTask := db.DB.Create(&task)
	createErr := createdTask.Error

	if createErr != nil {
		w.WriteHeader(http.StatusBadRequest) // 400
		w.Write([]byte(createErr.Error()))
		return
	}

	json.NewEncoder(w).Encode(&task)
	w.WriteHeader(http.StatusCreated)
}

func PatchTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	params := mux.Vars(r)
	db.DB.First(&task, params["id"])

	if task.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found :O"))
		return
	}

	// Mapping to filter PATCH fields requested
	var updates map[string]any
	jsonErr := json.NewDecoder(r.Body).Decode(&updates)
	if jsonErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid JSON"))
		return
	}

	// Delete fields that should not be updated
	delete(updates, "id")
	delete(updates, "created_at")
	delete(updates, "updated_at")

	updatedTask := db.DB.Model(&task).Updates(updates).First(&task)
	updateErr := updatedTask.Error
	if updateErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(updateErr.Error()))
		return
	}

	json.NewEncoder(w).Encode(&task)
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	params := mux.Vars(r)
	db.DB.First(&task, params["id"])

	if task.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Task not found :O"))
		return
	}

	db.DB.Unscoped().Delete(&task)
	w.WriteHeader(http.StatusNoContent) // 204
}
