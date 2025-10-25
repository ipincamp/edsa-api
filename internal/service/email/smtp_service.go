package email

import (
	"context"
	"fmt"
	"net/smtp"
	"strconv"

	"github.com/ipincamp/go-edsa-api/internal/config"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
)

type smtpService struct {
	cfg *config.Config
}

// NewSmtpService membuat instance EmailService baru
func NewSmtpService(cfg *config.Config) usecase.EmailService {
	return &smtpService{cfg: cfg}
}

// SendEmail mengirim email menggunakan konfigurasi SMTP dari .env
func (s *smtpService) SendEmail(ctx context.Context, to, subject, body string) error {
	cfg := s.cfg.Email

	// Setup autentikasi
	auth := smtp.PlainAuth("", cfg.User, cfg.Password, cfg.Host)

	// Alamat server SMTP
	addr := cfg.Host + ":" + strconv.Itoa(cfg.Port)

	// Format pesan email dengan header MIME untuk HTML
	msg := []byte(fmt.Sprintf(
		"To: %s\r\n"+
			"From: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-version: 1.0;\n"+
			"Content-Type: text/html; charset=\"UTF-8\";\n\n"+
			"%s\r\n",
		to, cfg.SenderEmail, subject, body,
	))

	// Kirim email
	return smtp.SendMail(addr, auth, cfg.SenderEmail, []string{to}, msg)
}
