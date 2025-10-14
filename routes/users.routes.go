package routes

import (
	"encoding/json"
	"net/http"

	"github.com/Cruz000Testing/go-gorm-restapi/db"
	"github.com/Cruz000Testing/go-gorm-restapi/models"
	"github.com/gorilla/mux"
)

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	db.DB.Find(&users)
	json.NewEncoder(w).Encode(&users)
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	params := mux.Vars(r)
	db.DB.First(&user, params["id"])

	if user.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found :O"))
		return
	}

	db.DB.Model(&user).Association("Task").Find(&user.Tasks)

	json.NewEncoder(w).Encode(&user)
}

func PostUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	jsonErr := json.NewDecoder(r.Body).Decode(&user)

	if jsonErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid JSON"))
		return
	}

	createdUser := db.DB.Create(&user)
	createErr := createdUser.Error

	if createErr != nil {
		w.WriteHeader(http.StatusBadRequest) // 400
		w.Write([]byte(createErr.Error()))
		return
	}

	json.NewEncoder(w).Encode(&user)
	w.WriteHeader(http.StatusCreated)
}

func PatchUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	params := mux.Vars(r)
	db.DB.First(&user, params["id"])

	if user.ID == 0 {
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

	updatedUser := db.DB.Model(&user).Updates(updates).First(&user)
	updateErr := updatedUser.Error

	if updateErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(updateErr.Error()))
		return
	}

	json.NewEncoder(w).Encode(&user)
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	params := mux.Vars(r)
	db.DB.First(&user, params["id"])

	if user.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found :O"))
		return
	}

	db.DB.Unscoped().Delete(&user)
	w.WriteHeader(http.StatusNoContent) // 204
}
