package utils

import (
	"fmt"
	"log"
	"net/smtp"
	"net/http"
	"os"
	"encoding/json"
)

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func SendEmail(to, subject, body string) error {
	from := os.Getenv("SMTP_EMAIL")
	pass := os.Getenv("SMTP_PASSWORD")
	hostname := os.Getenv("SMTP_HOSTNAME")

	if from == "" || pass == "" {
		return fmt.Errorf("SMTP credentials are not set")
	}

	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", from, to, subject, body)
	err := smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", from, pass,hostname), from, []string{to}, []byte(msg))
	if err != nil {
		log.Printf("Failed to send email to %s: %v", to, err)
		return err
	}

	log.Printf("Email sent successfully to %s", to)
	return nil
}
