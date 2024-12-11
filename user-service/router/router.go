package router

import (
	"net/http"
	"github.com/gorilla/mux"
	"user-service/handlers"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Task service is running"))
	}).Methods(http.MethodGet)
	r.HandleFunc("/signup", handlers.SignUpHandler).Methods("POST")
	r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	return r
}
