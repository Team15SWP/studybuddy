package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"study_buddy/internal/config"
	"study_buddy/internal/model"

	"gopkg.in/gomail.v2"
)

var _ Service = (*NotifyService)(nil)

type NotifyService struct {
	repo       NotifyProvider
	log        *slog.Logger
	promptsCfg *config.Prompts
	smtpCfg    *config.SmtpConfig
	notifyCfg  *config.NotifyConfig
	openAiCfg  *config.OpenAI
}

func NewNotifyService(repo NotifyProvider, log *slog.Logger, cfg *config.Config) *NotifyService {
	return &NotifyService{
		repo:       repo,
		log:        log,
		promptsCfg: &cfg.Prompts,
		smtpCfg:    &cfg.SmtpConfig,
		notifyCfg:  &cfg.NotifyConfig,
		openAiCfg:  &cfg.OpenAI,
	}
}

type Service interface {
	NotifyUsers(ctx context.Context, now *time.Time) error
	SendMessage(ctx context.Context, messageBody, toEmail, name string) error
}

type NotifyProvider interface {
	GetAllUsersEmail(ctx context.Context, userIDs []int64) ([]*model.User, error)
	GetUserIDs(ctx context.Context, now *time.Time) ([]int64, error)
}

func (n *NotifyService) NotifyUsers(ctx context.Context, now *time.Time) error {
	n.log.Info("NotifyUsers start")
	userIDs, err := n.repo.GetUserIDs(ctx, now)
	n.log.Info("userIDs", len(userIDs), userIDs)
	users, err := n.repo.GetAllUsersEmail(ctx, userIDs)
	if err != nil {
		return fmt.Errorf("n.repo.GetAllUsersEmail: %w", err)
	}
	for _, user := range users {
		n.log.Info("user info", user.Name, user.Email)
	}
	messageBody := `
<p>Hi, %s!</p>
<p>Welcome to Study Buddy — your daily place to grow with code 😊</p>
<p>We’ve got fresh programming challenges waiting for you every day — come back often, stay sharp, and have fun!</p>
<p>See you inside,<br>The Study Buddy Team</p>
`
	for _, user := range users {
		err = n.SendMessage(ctx, messageBody, user.Email, user.Name)
		if err != nil {
			n.log.Error(fmt.Sprintf("n.SendMessage: %v", err))
		}
	}
	return nil
}

func (n *NotifyService) SendMessage(ctx context.Context, messageBody, toEmail, name string) error {
	body := fmt.Sprintf(messageBody, name)

	m := gomail.NewMessage()

	m.SetHeader("From", n.notifyCfg.From)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", n.notifyCfg.Subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(n.smtpCfg.Host, n.smtpCfg.Port, n.smtpCfg.Email, n.smtpCfg.Password)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("gomail.NewDialer: %w", err)
	}
	return nil
}
