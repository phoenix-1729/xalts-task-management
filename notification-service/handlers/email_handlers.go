package handlers

import (
	"encoding/json"
	"net/http"
	"notification-service/utils"
)

func TaskCreatedHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		UserID int    `json:"user_id"`
		Title  string `json:"title"`
		UserEmail string `json:"user_email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	email := payload.UserEmail
	if email == "" {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch user email")
		return
	}

	subject := "New Task Assigned for Approval"
	body := "You have been assigned to approve the task: " + payload.Title
	if err := utils.SendEmail(email, subject, body); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Notification sent"})
}

func TaskApprovedHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		UserID int    `json:"user_id"`
		TaskID int    `json:"task_id"`
		Status string `json:"status"`
		Comment string `json:"comment"`
		UserEmail string `json:"user_email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	email := payload.UserEmail
	if email == "" {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch user email")
		return
	}

	subject := "Task Approval Update"
	body := "Your approval status for Task ID " + string(payload.TaskID) + " is: " + payload.Status + ". Comment: " + payload.Comment
	if err := utils.SendEmail(email, subject, body); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Notification sent"})
}
