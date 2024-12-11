package main

import (
	"log"
	"net/http"
	"task-service/database"
	"task-service/router"
)

func main() {
	database.InitDB()
	defer database.CloseDB()

	r := router.InitializeRouter()
	log.Println("Task service is running on port 8001")
	log.Fatal(http.ListenAndServe(":8001", r))
}
