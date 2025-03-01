package email

import (
	"crypto/tls"
	"strconv"

	"github.com/hewpao/hewpao-backend/config"
	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/dto"
	"github.com/hewpao/hewpao-backend/repository"
	"gopkg.in/gomail.v2"
)

type gmailEmailNotificationRepo struct {
	message  *gomail.Message
	cfg      config.Config
	gmailCli *gomail.Dialer
}

func NewGmailEmailNotificationRepo(message *gomail.Message, cfg *config.Config) (repository.NotificationRepository, error) {
	emailPort, err := strconv.Atoi(cfg.EmailPort)
	if err != nil {
		return nil, err
	}

	gmailCli := gomail.NewDialer(cfg.EmailServer, emailPort, cfg.EmailUser, cfg.EmailPassword)
	gmailCli.SSL = true
	gmailCli.TLSConfig = &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         cfg.EmailServer,
	}

	return &gmailEmailNotificationRepo{
		message:  message,
		cfg:      *cfg,
		gmailCli: gmailCli,
	}, nil
}

func (e *gmailEmailNotificationRepo) Notify(toUser *domain.User, req *dto.NotificationDTO) error {
	e.message.SetHeader("From", e.cfg.EmailUser)
	e.message.SetHeader("To", toUser.Email)
	e.message.SetHeader("Subject", req.Subject)
	e.message.SetBody("text/html", req.Content)

	err := e.gmailCli.DialAndSend(e.message)
	if err != nil {
		return err
	}

	return nil
}
