package email

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"sync"
)

func SendTaskCreatedNotification(to string, taskTitle string) error {
	subject := "New Task Assigned for Approval"
	body := fmt.Sprintf("You have been assigned to approve the task: %s", taskTitle)
	return sendEmail(to, subject, body)
}

func NotifyTaskCreatorOnSignOff(creatorEmail string, approverName string, taskTitle string) error {
	subject := "Task Approval Update"
	body := fmt.Sprintf("%s has signed off on the task: %s.", approverName, taskTitle)
	return sendEmail(creatorEmail, subject, body)
}

func NotifyAllPartiesOnCompletion(emails []string, taskTitle string) error {
	subject := "Task Fully Approved"
	body := fmt.Sprintf("The task: %s has been fully approved by all approvers.", taskTitle)

	var wg sync.WaitGroup
	errChan := make(chan error, len(emails))

	for _, email := range emails {
		wg.Add(1)
		go func(email string) {
			defer wg.Done()
			if err := sendEmail(email, subject, body); err != nil {
				errChan <- fmt.Errorf("failed to send email to %s: %v", email, err)
			}
		}(email)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		log.Println("Error:", err)
		return err
	}

	return nil
}

func sendEmail(to, subject, body string) error {
	from := os.Getenv("SMTP_EMAIL")
	pass := os.Getenv("SMTP_PASSWORD")

	if from == "" || pass == "" {
		return fmt.Errorf("SMTP credentials are not set")
	}

	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", from, to, subject, body)
	err := smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", from, pass), from, []string{to}, []byte(msg))
	if err != nil {
		return fmt.Errorf("failed to send email to %s: %w", to, err)
	}

	log.Printf("Email sent successfully to %s\n", to)
	return nil
}
