package main

import (
	"log"
	"net/http"
	"notification-service/router"
)

func main() {
	r := router.SetupRouter()
	log.Println("Notification service is running on port 8002")
	log.Fatal(http.ListenAndServe(":8002", r))
}
