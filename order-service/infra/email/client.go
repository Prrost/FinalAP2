// infra/email/client.go
package email

import (
	"fmt"
	"net/smtp"
)

type EmailSender struct {
	Host     string
	Port     string
	Username string
	Password string
}

// NewEmailSenderConfig создаёт клиента с явными параметрами
func NewEmailSenderConfig(host, port, username, password string) *EmailSender {
	return &EmailSender{Host: host, Port: port, Username: username, Password: password}
}

func (e *EmailSender) Send(to, subject, body string) error {
	auth := smtp.PlainAuth("", e.Username, e.Password, e.Host)
	msg := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body))
	addr := fmt.Sprintf("%s:%s", e.Host, e.Port)
	return smtp.SendMail(addr, auth, e.Username, []string{to}, msg)
}
