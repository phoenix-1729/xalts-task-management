package router

import (
	"github.com/gorilla/mux"
	"net/http"
	"notification-service/utils"
	"notification-service/handlers"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		utils.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
	}).Methods("GET")

	r.HandleFunc("/notify/task-created", handlers.TaskCreatedHandler).Methods("POST")
	r.HandleFunc("/notify/task-approved", handlers.TaskApprovedHandler).Methods("POST")
	r.HandleFunc("/notify/task-rejected", handlers.TaskApprovedHandler).Methods("POST")

	return r
}
