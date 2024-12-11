package main

import (
	"log"
	"net/http"
	"user-service/database"
	"user-service/router"
)

func main() {
	database.InitDB()
	defer database.CloseDB()

	r := router.SetupRouter()
	log.Println("User service is running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
