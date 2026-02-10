package email

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// IsSimulationMode returns true when SIMULATE_LOGIN_EMAILS=true (localhost testing without SMTP).
func IsSimulationMode() bool {
	return os.Getenv("SIMULATE_LOGIN_EMAILS") == "true"
}

// ShouldSendLoginEmails returns true when we should send (or simulate) login emails:
// - On Render with SMTP configured, or SEND_LOGIN_EMAILS with SMTP, or SIMULATE_LOGIN_EMAILS (no SMTP needed).
func ShouldSendLoginEmails() bool {
	if os.Getenv("RENDER") == "true" || os.Getenv("SEND_LOGIN_EMAILS") == "true" || IsSimulationMode() {
		return IsSimulationMode() || smtpConfigured()
	}
	return false
}

func smtpConfigured() bool {
	host := os.Getenv("SMTP_HOST")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASSWORD")
	return host != "" && user != "" && pass != ""
}

// SimulatedEmailPath is where the last simulated email is written (when SIMULATE_LOGIN_EMAILS=true).
func SimulatedEmailPath() string {
	p := os.Getenv("SIMULATE_LOGIN_EMAIL_FILE")
	if p == "" {
		p = "data/simulated_login_email.txt"
	}
	return p
}

// SendLoginNotification sends an email to the given address when a reader or consultant logs in.
// When SIMULATE_LOGIN_EMAILS=true, writes the would-be email to a file and logs instead of sending.
func SendLoginNotification(to, role, name, userEmail string) {
	subject := fmt.Sprintf("[Alice Suite] %s logged in", role)
	body := fmt.Sprintf("A %s has logged in.\n\nName: %s\nEmail: %s\n\nThis is an automated notification from Alice Suite.", role, name, userEmail)

	if IsSimulationMode() {
		path := SimulatedEmailPath()
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			log.Printf("email: simulate: could not create dir for %s: %v", path, err)
			return
		}
		content := fmt.Sprintf("To: %s\nSubject: %s\n\n%s", to, subject, body)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			log.Printf("email: simulate: could not write %s: %v", path, err)
			return
		}
		log.Printf("email: simulated login notification written to %s (To: %s, %s %s)", path, to, role, userEmail)
		return
	}

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
