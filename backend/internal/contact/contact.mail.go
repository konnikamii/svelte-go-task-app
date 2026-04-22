package contact

import (
	"bytes"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type MailConfig struct {
	Host string
	Port int
	From string
	To   string
}

func (cfg MailConfig) SendContactNotification(request ContactRequestResponse) error {
	host := strings.TrimSpace(cfg.Host)
	from := strings.TrimSpace(cfg.From)
	to := strings.TrimSpace(cfg.To)
	if host == "" || cfg.Port <= 0 || from == "" || to == "" {
		return fmt.Errorf("contact notification mail config is incomplete")
	}

	addr := fmt.Sprintf("%s:%d", host, cfg.Port)
	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := conn.SetDeadline(time.Now().Add(5 * time.Second)); err != nil {
		return err
	}

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	defer client.Close()

	if err := client.Mail(from); err != nil {
		return err
	}
	if err := client.Rcpt(to); err != nil {
		return err
	}

	writer, err := client.Data()
	if err != nil {
		return err
	}

	var message bytes.Buffer
	message.WriteString(fmt.Sprintf("To: %s\r\n", to))
	message.WriteString(fmt.Sprintf("From: %s\r\n", from))
	message.WriteString(fmt.Sprintf("Reply-To: %s\r\n", request.Email))
	message.WriteString(fmt.Sprintf("Subject: New contact request: %s\r\n", request.Title))
	message.WriteString("MIME-Version: 1.0\r\n")
	message.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	message.WriteString("\r\n")
	message.WriteString("A new contact request was submitted.\r\n\r\n")
	message.WriteString(fmt.Sprintf("ID: %d\r\n", request.ID))
	message.WriteString(fmt.Sprintf("Email: %s\r\n", request.Email))
	message.WriteString(fmt.Sprintf("Title: %s\r\n\r\n", request.Title))
	message.WriteString("Message:\r\n")
	message.WriteString(request.Message)
	message.WriteString("\r\n")

	if _, err := writer.Write(message.Bytes()); err != nil {
		_ = writer.Close()
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}

	return client.Quit()
}

func logContactMailFailure(err error, request ContactRequestResponse) {
	logrus.WithError(err).WithFields(logrus.Fields{
		"contact_request_id": request.ID,
		"contact_email":      request.Email,
	}).Warn("contact request saved but notification email was not sent")
}
