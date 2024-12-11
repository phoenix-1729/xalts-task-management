package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
	"user-service/models"
	"user-service/utils"

	"github.com/golang-jwt/jwt/v4"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("LoginHandler invoked")
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("Error decoding request body: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	log.Printf("Attempting login for user: %s", user.Email)
	storedUser, err := models.GetUserByEmail(user.Email)
	if err != nil {
		log.Printf("Error fetching user by email: %v", err)
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	if storedUser.Password != user.Password {
		log.Printf("Password mismatch for user: %s", user.Email)
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	log.Printf("Generating token for user ID: %d", storedUser.ID)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": storedUser.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Printf("Error generating token: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	log.Printf("Token generated successfully for user ID: %d", storedUser.ID)
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"token": tokenString})
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("SignUpHandler invoked")
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("Error decoding request body: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	log.Printf("Attempting to create user: %s", user.Email)
	id, err := models.CreateUser(user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("User created successfully with ID: %d", id)
	utils.RespondWithJSON(w, http.StatusCreated, map[string]interface{}{"id": id})
}
