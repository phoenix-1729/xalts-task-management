package models

import (
	"user-service/database"
	"log"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateUser(user User) (int, error) {
	db := database.GetDB()

	var id int
	err := db.QueryRow("INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id", user.Name, user.Email, user.Password).Scan(&id)
	if err != nil {
		log.Println("Error inserting user:", err)
		return 0, err
	}

	return id, nil
}

func GetUserByEmail(email string) (User, error) {
	db := database.GetDB()

	var user User
	err := db.QueryRow("SELECT id, name, email, password FROM users WHERE email = $1", email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		log.Println("Error fetching user by email:", err)
		return user, err
	}

	return user, nil
}