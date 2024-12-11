package router

import (
	"net/http"
	"task-service/handlers"
	"task-service/middleware"

	"github.com/gorilla/mux"
)

func InitializeRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Task service is running"))
	}).Methods(http.MethodGet)

	protected := router.PathPrefix("/").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.HandleFunc("/create-task", handlers.CreateTaskHandler).Methods(http.MethodPost)
	protected.HandleFunc("/tasks/approve", handlers.ApproveTaskHandler).Methods(http.MethodPost)
	protected.HandleFunc("/tasks", handlers.ViewTaskStatusHandler).Methods(http.MethodGet)
	return router
}
