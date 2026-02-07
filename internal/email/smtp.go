package email

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strconv"
	"strings"
)

// ShouldSendLoginEmails returns true when running on Render (or SEND_LOGIN_EMAILS=true) and SMTP is configured.
// Call this before sending so we never send on localhost unless explicitly enabled.
func ShouldSendLoginEmails() bool {
	if os.Getenv("RENDER") != "true" && os.Getenv("SEND_LOGIN_EMAILS") != "true" {
		return false
	}
	return smtpConfigured()
}

func smtpConfigured() bool {
	host := os.Getenv("SMTP_HOST")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASSWORD")
	return host != "" && user != "" && pass != ""
}

// SendLoginNotification sends an email to the given address when a reader or consultant logs in.
// No-op if SMTP is not configured (logs and returns nil).
func SendLoginNotification(to, role, name, userEmail string) {
	if !smtpConfigured() {
		log.Printf("email: SMTP not configured, skipping login notification to %s", to)
		return
	}
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	if port == "" {
		port = "587"
	}
	portNum, _ := strconv.Atoi(port)
	from := os.Getenv("SMTP_FROM")
	if from == "" {
		from = os.Getenv("SMTP_USER")
	}
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASSWORD")

	subject := fmt.Sprintf("[Alice Suite] %s logged in", role)
	body := fmt.Sprintf("A %s has logged in.\n\nName: %s\nEmail: %s\n\nThis is an automated notification from Alice Suite.", role, name, userEmail)
	msg := "From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body + "\r\n"

	addr := fmt.Sprintf("%s:%d", host, portNum)
	auth := smtp.PlainAuth("", user, pass, strings.TrimSuffix(host, ":25"))
	err := smtp.SendMail(addr, auth, from, []string{to}, []byte(msg))
	if err != nil {
		log.Printf("email: failed to send login notification: %v", err)
		return
	}
	log.Printf("email: login notification sent to %s for %s %s", to, role, userEmail)
}
