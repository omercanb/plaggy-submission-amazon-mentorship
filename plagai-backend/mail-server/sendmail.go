package mailserver

import (
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

// MailServer holds SMTP configuration
type MailServer struct {
	Host     string
	Port     string
	User     string
	Password string
}

// MailMessage holds email details
type MailMessage struct {
	FromEmail string
	FromName  string
	ToEmail   string
	Subject   string
	Content   string
}

func SendMail(server MailServer, message MailMessage) error {
	msg := []byte(
		fmt.Sprintf("From: %s <%s>\r\n", message.FromName, message.FromEmail) +
			fmt.Sprintf("To: %s\r\n", message.ToEmail) +
			fmt.Sprintf("Subject: %s\r\n", message.Subject) +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=\"utf-8\"\r\n\r\n" +
			message.Content + "\r\n",
	)

	auth := smtp.PlainAuth("", server.User, server.Password, server.Host)
	return smtp.SendMail(server.Host+":"+server.Port, auth, message.FromEmail, []string{message.ToEmail}, msg)
}

func SendMailWithPlaintextContent(to string, fromname string, subject string, content string) {
	godotenv.Load("../.env")
	pass := os.Getenv("BREVO_SMTP_KEY")
	if pass == "" {
		log.Fatal("BREVO_SMTP_KEY not found")
	}

	server := MailServer{
		Host:     "smtp-relay.brevo.com",
		Port:     "587",
		User:     "9607cd001@smtp-brevo.com",
		Password: pass,
	}

	message := MailMessage{
		FromEmail: "omercanbaykara@gmail.com",
		FromName:  fromname,
		ToEmail:   to,
		Subject:   subject,
		Content:   content,
	}

	if err := SendMail(server, message); err != nil {
		log.Fatal(err)
	}

	log.Println("Email sent successfully")
}
