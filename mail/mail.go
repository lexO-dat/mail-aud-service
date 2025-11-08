package mail

import (
	"fmt"
	"log"
	"net/smtp"
	"time"

	"github.com/jordan-wright/email"
)

const (
	smtpAuthAdress = "smtp.gmail.com"
	smtpServerAuth = "smtp.gmail.com:587"
)

type EmailSender interface {
	SendEmail(
		subject string,
		body string,
		to []string,
		cc []string,
		bcc []string,
		attachFiles []string,
	) error
}

type GmailSender struct {
	name              string
	fromEmailAdress   string
	fromEmailPassword string
}

func NewGmailSender(name string, fromEmailAdress string, fromEmailPassword string) EmailSender {
	return &GmailSender{
		name:              name,
		fromEmailAdress:   fromEmailAdress,
		fromEmailPassword: fromEmailPassword,
	}
}

func (sender *GmailSender) SendEmail(
	subject string,
	body string,
	to []string,
	cc []string,
	bcc []string,
	attachFiles []string,
) error {
	log.Printf("📧 Iniciando envío de email a: %v", to)
	
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAdress)
	e.Subject = subject
	e.HTML = []byte(body)
	e.To = to
	e.Cc = cc
	e.Bcc = bcc

	log.Printf("📎 Adjuntando %d archivos...", len(attachFiles))
	for _, f := range attachFiles {
		_, err := e.AttachFile(f)
		if err != nil {
			return fmt.Errorf("error attaching file: %s", err)
		}
	}

	log.Println("🔐 Configurando autenticación SMTP...")
	smptpAuth := smtp.PlainAuth("", sender.fromEmailAdress, sender.fromEmailPassword, smtpAuthAdress)

	log.Printf("📡 Conectando a servidor SMTP: %s", smtpServerAuth)
	
	// Crear un canal para manejar el timeout
	done := make(chan error, 1)
	
	go func() {
		done <- e.Send(smtpServerAuth, smptpAuth)
	}()
	
	// Esperar por el resultado o timeout
	select {
	case err := <-done:
		if err != nil {
			log.Printf("❌ Error SMTP: %v", err)
			return err
		}
		log.Println("✅ Email enviado exitosamente por SMTP")
		return nil
	case <-time.After(30 * time.Second):
		log.Println("⏰ Timeout al enviar email (30s)")
		return fmt.Errorf("timeout sending email after 30 seconds")
	}
}
