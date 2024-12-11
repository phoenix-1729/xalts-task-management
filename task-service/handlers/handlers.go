package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"task-service/models"
	"task-service/utils"
	"fmt"
)

const NotificationServiceURL = "http://localhost:8082"

func ViewTaskStatusHandler(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("id")
	if taskID == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Task ID is required")
		return
	}

	task, err := models.GetTaskByID(taskID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error retrieving task: %v", err))
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, task)
}
func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	taskID, err := models.CreateTask(task)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	for _, approverID := range task.Approvers {
		email, err := models.GetUserEmailByID(approverID)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		notifyPayload := map[string]interface{}{
			"user_id": approverID,
			"title":   task.Title,
			"user_email": email,
		}
		go sendNotification("/notify/task-created", notifyPayload)
	}

	utils.RespondWithJSON(w, http.StatusCreated, map[string]interface{}{"task_id": taskID})
}

func ApproveTaskHandler(w http.ResponseWriter, r *http.Request) {
	var approval models.Approval
	if err := json.NewDecoder(r.Body).Decode(&approval); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	message, err := models.ApproveTask(approval)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	email, err := models.GetUserEmailByID(approval.UserID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	notifyPayload := map[string]interface{}{
		"user_id": approval.UserID,
		"task_id": approval.TaskID,
		"comment": approval.Comment,
		"user_email": email,
	}
	if approval.IsSigned {
		notifyPayload["status"] = "Signed"
		go sendNotification("/notify/task-approved", notifyPayload)
	} else {
		notifyPayload["status"] = "Rejected"
		go sendNotification("/notify/task-rejected", notifyPayload)
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": message})
}

func sendNotification(endpoint string, payload interface{}) {
	data, _ := json.Marshal(payload)
	_, _ = http.Post(NotificationServiceURL+endpoint, "application/json", bytes.NewBuffer(data))
}
