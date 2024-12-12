package auth

import (
	"fmt"
	"math/rand"
	"net/smtp"
	"os"
	"time"
)

type EmailManager struct {
	host     string
	port     string
	password string
	email    string
}

func NewEmailManager() *EmailManager {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	password := os.Getenv("SMTP_PASSWORD")
	email := os.Getenv("FROM_EMAIL")
	emailManager := EmailManager{
		host:     host,
		port:     port,
		password: password,
		email:    email,
	}
	if host == "" || port == "" || password == "" || email == "" {
		return nil
	}
	return &emailManager
}
func (em *EmailManager) GenerateVerifyCode() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06d", r.Intn(1000000))
}
func (em *EmailManager) SendCode(username string, userEmail string, verifyCode string) error {
	auth := smtp.PlainAuth(
		"",
		em.email,
		em.password,
		em.host,
	)
	toEmail := userEmail
	msg := fmt.Sprintf("Hello, %s!\nVerification code for password recovery: %s.", username, verifyCode)
	htmlBody := fmt.Sprintf("To: %s\nSubject:Recovery\n%s", toEmail, msg)

	err := smtp.SendMail(
		fmt.Sprintf("%s:%s", em.host, em.port),
		auth,
		em.email,
		[]string{toEmail},
		[]byte(htmlBody),
	)
	if err != nil {
		return fmt.Errorf("unable to send a mail: %s", err)
	}

	return nil
}
