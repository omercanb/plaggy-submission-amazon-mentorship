package mailserver

import (
	"log"
	"net/smtp"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMailSending(t *testing.T) {
	godotenv.Load("../.env")
	pass := os.Getenv("BREVO_SMTP_KEY")
	if pass == "" {
		t.Fatal("BREVO_SMTP_KEY not found")
	}

	from := "omercanbaykara@gmail.com"
	to := "omercanbaykara@gmail.com"
	msg := []byte(
		"From: Login at PlagAI <omercanbaykara@gmail.com>\r\n" +
			"To: omercanbaykara@gmail.com\r\n" +
			"Subject: Test via Brevo SMTP\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=\"utf-8\"\r\n\r\n" +
			"Hello from Go + Brevo!\r\n",
	)
	// Brevo SMTP config
	smtpHost := "smtp-relay.brevo.com"
	smtpPort := "587"
	auth := smtp.PlainAuth("", "9607cd001@smtp-brevo.com", pass, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Email sent successfully!")
}
