package notitype

import (
	"crypto/tls"
	"html/template"
	"strconv"
	"strings"
	"time"

	"github.com/hewpao/hewpao-backend/config"
	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/dto"
	"github.com/hewpao/hewpao-backend/repository"
	"gopkg.in/gomail.v2"
)

type emailNotificationRepo struct {
	message  *gomail.Message
	cfg      *config.Config
	gmailCli *gomail.Dialer
}

func NewEmailNotificationRepo(message *gomail.Message, cfg *config.Config) (repository.NotificationRepository, error) {
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

	return &emailNotificationRepo{
		message:  message,
		cfg:      cfg,
		gmailCli: gmailCli,
	}, nil
}

func prContentBuilder(productRequest *domain.ProductRequest, cfg *config.Config, content *strings.Builder) error {
	data := dto.NotificationDataDTO{
		RecipientName: productRequest.User.Name,
		CompanyName:   "HEWPAO",
		ProductID:     productRequest.ID,
		ProductStatus: productRequest.DeliveryStatus,
		SupportEmail:  cfg.EmailUser,
		Year:          time.Now().Year(),
	}

	tmpl, err := template.ParseFiles("./assets/productRequest/emailTemplate.html")
	if err != nil {
		return err
	}

	content.Reset()

	err = tmpl.Execute(content, data)
	if err != nil {
		return err
	}

	return nil
}

func (e *emailNotificationRepo) PrNotify(user *domain.User, prod *domain.ProductRequest) error {
	var content strings.Builder
	err := prContentBuilder(prod, e.cfg, &content)
	if err != nil {
		return err
	}

	e.message.SetHeader("From", e.cfg.EmailUser)
	e.message.SetHeader("To", user.Email)
	e.message.SetHeader("Subject", "[HEWPAO] Product Request Status Notification!")
	e.message.SetBody("text/html", content.String())

	err = e.gmailCli.DialAndSend(e.message)
	if err != nil {
		return err
	}

	return nil
}
