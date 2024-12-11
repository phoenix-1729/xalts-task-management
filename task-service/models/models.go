package models

import (
	"log"
	"task-service/database"

	"github.com/lib/pq"
)

type Task struct {
	ID        int      `json:"id"`
	Title     string   `json:"title"`
	Status    string   `json:"status"`
	CreatorID int      `json:"creator_id"`
	Approvers []int    `json:"approvers"`
	Comments  []string `json:"comments"`
}

type Approval struct {
	TaskID    int    `json:"task_id"`
	UserID    int    `json:"user_id"`
	Comment   string `json:"comment"`
	IsSigned  bool   `json:"is_signed"`
}

func CreateTask(task Task) (int, error) {
	log.Println("Creating new task...")
	db := database.GetDB()

	query := `
		INSERT INTO tasks (title, status, creator_id, approvers)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var taskID int
	err := db.QueryRow(query, task.Title, "pending", task.CreatorID, pq.Array(task.Approvers)).Scan(&taskID)
	if err != nil {
		log.Printf("Error inserting task: %v", err)
		return 0, err
	}
	log.Printf("Task created successfully with ID: %d", taskID)

	for _, approver := range task.Approvers {
		query = `
			INSERT INTO approvals (task_id, user_id, is_signed)
			VALUES ($1, $2, $3)
		`
		_, err := db.Exec(query, taskID, approver, false)
		if err != nil {
			log.Printf("Error inserting approver (Task ID: %d, User ID: %d): %v", taskID, approver, err)
			return 0, err
		}
		log.Printf("Approver (User ID: %d) added to task (Task ID: %d)", approver, taskID)
	}

	return taskID, nil
}

func ApproveTask(approval Approval) (string, error) {
	log.Printf("Recording approval for Task ID: %d by User ID: %d", approval.TaskID, approval.UserID)
	db := database.GetDB()

	query := `
		UPDATE approvals
		SET is_signed = $1
		WHERE task_id = $2 AND user_id = $3
	`
	_, err := db.Exec(query, true, approval.TaskID, approval.UserID)
	if err != nil {
		log.Printf("Error updating approval: %v", err)
		return "", err
	}
	log.Printf("Approval recorded for Task ID: %d by User ID: %d", approval.TaskID, approval.UserID)

	var pendingApprovals int
	query = `
		SELECT COUNT(*)
		FROM approvals
		WHERE task_id = $1 AND is_signed = $2
	`
	err = db.QueryRow(query, approval.TaskID, false).Scan(&pendingApprovals)
	if err != nil {
		log.Printf("Error checking pending approvals: %v", err)
		return "", err
	}

	if pendingApprovals == 0 {
		query = `
			UPDATE tasks
			SET status = $1
			WHERE id = $2
		`
		_, err := db.Exec(query, "approved", approval.TaskID)
		if err != nil {
			log.Printf("Error updating task status to approved: %v", err)
			return "", err
		}
		log.Printf("Task ID: %d approved by all approvers.", approval.TaskID)
		return "Task approved by all users.", nil
	}

	log.Printf("Approval recorded for Task ID: %d. Awaiting other approvers.", approval.TaskID)
	return "Approval recorded. Awaiting other users.", nil
}

func GetUserEmailByID(userID int) (string, error) {
    db := database.GetDB()

    var email string
    err := db.QueryRow("SELECT email FROM users WHERE id = $1", userID).Scan(&email)
    if err != nil {
        log.Println("Error fetching email by userID:", err)
        return "", err
    }

    return email, nil
}

func GetTaskByID(taskID string) (*Task, error) {
	var task Task
	db := database.GetDB()
	query := `SELECT id, title, creator_id, status 
	          FROM tasks WHERE id = $1`

	err := db.QueryRow(query, taskID).Scan(&task.ID, &task.Title, &task.CreatorID, &task.Status)
	if err != nil {
		return nil, err
	}

	return &task, nil
}