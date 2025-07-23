package service

import (
	"context"
	"fmt"
	"time"

	"study_buddy/internal/config"
	"study_buddy/internal/model"
	"study_buddy/pkg/constants"
	"study_buddy/pkg/errlist"
	"study_buddy/pkg/hash"

	"github.com/dgrijalva/jwt-go"
	"gopkg.in/gomail.v2"
)

var _ Service = (*AuthService)(nil)

type AuthService struct {
	repo             UserProvider
	notificationRepo NotificationProvider
	hashConfig       *config.HashConfig
	smtpCfg          *config.SmtpConfig
	notifyCfg        *config.NotifyConfig
}

func NewAuthService(repo UserProvider, notificationRepo NotificationProvider, cfg *config.Config) *AuthService {
	return &AuthService{
		repo:             repo,
		notificationRepo: notificationRepo,
		hashConfig:       &cfg.HashConfig,
		smtpCfg:          &cfg.SmtpConfig,
		notifyCfg:        &cfg.NotifyConfig,
	}
}

type Service interface {
	LogIn(ctx context.Context, email, password string) (*model.AuthToken, error)
	SignUp(ctx context.Context, username, email, password string) (*VerifyAccount, error)
	SendConfirmationMessage(toEmail, token string) error
	Confirm(ctx context.Context, token string) error
}

type UserProvider interface {
	GetUserByEmailOrUsername(ctx context.Context, username string) (*model.UserData, error)
	CreateUser(ctx context.Context, username, email, password string) (*model.UserData, error)
	UpdateUser(ctx context.Context, user *model.UserData) error
}

type NotificationProvider interface {
	CreateNotification(ctx context.Context, notif *model.Notification) error
}

func (a *AuthService) LogIn(ctx context.Context, username, password string) (*model.AuthToken, error) {
	user, err := a.repo.GetUserByEmailOrUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("[authService][LogIn][GetUserByEmailOrUsername]: %w", err)
	}
	if !user.IsConfirmed {
		return nil, fmt.Errorf("[authService][LogIn][IsConfirmed]: %w", errlist.ErrUserIsNotVerified)
	}
	if !hash.ComparePassword(password, user.Password) {
		return nil, fmt.Errorf("[authService][LogIn][ComparePassword]: %w", errlist.ErrPasswordIsIncorrect)
	}

	tokenString, err := signToken(user, a.hashConfig.SigningKey)
	if err != nil {
		return nil, fmt.Errorf("[authService][LogIn][SignToken]: %w", err)
	}
	return model.NewAuthToken(tokenString, user.Role), nil
}

type VerifyAccount struct {
	Message string `json:"message"`
}

func (a *AuthService) SignUp(ctx context.Context, username, email, password string) (*VerifyAccount, error) {
	if _, err := a.repo.GetUserByEmailOrUsername(ctx, username); err == nil {
		return nil, fmt.Errorf("[authService][SignUp][GetUserByEmailOrUsername]: %w", errlist.ErrUserExists)
	}
	if _, err := a.repo.GetUserByEmailOrUsername(ctx, email); err == nil {
		return nil, fmt.Errorf("[authService][SignUp][GetUserByEmailOrUsername]: %w", errlist.ErrUserExists)
	}

	hashedPassword, err := hash.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("[authService][SignUp][HashPassword]: %w", err)
	}

	user, err := a.repo.CreateUser(ctx, username, email, hashedPassword)
	if err != nil {
		return nil, fmt.Errorf("[authService][SignUp][CreateUser]: %w", err)
	}

	token, err := generateConfirmationToken(user.ID, a.hashConfig.SigningKey)
	if err != nil {
		return nil, fmt.Errorf("[authService][SignUp][GenerateConfirmationToken]: %w", err)
	}
	if err = a.SendConfirmationMessage(email, token); err != nil {
		return nil, fmt.Errorf("[authService][SignUp][SendConfirmationMessage]: %w", err)
	}

	tt, err := time.Parse("15:04", "00:00")
	if err != nil {
		return nil, fmt.Errorf("[authService][SignUp][ParseTime]: %w", err)
	}

	err = a.notificationRepo.CreateNotification(ctx, &model.Notification{
		UserID:  user.ID,
		Enabled: false,
		Time24:  tt,
		Days:    make([]int, 0, 7),
	})
	if err != nil {
		return nil, fmt.Errorf("[authService][SignUp][CreateNotification]: %w", err)
	}
	return &VerifyAccount{Message: "check your email and verify account"}, nil
}

func signToken(user *model.UserData, signingKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		constants.UserID: user.ID,
		constants.Role:   user.Role,
		constants.Name:   user.Name,
		constants.Email:  user.Email,
	})

	secretKey := []byte(signingKey)
	tokenString, err := token.SignedString(secretKey)
	return tokenString, err
}

func generateConfirmationToken(userId int64, signingKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		constants.UserID: userId,
		constants.Exp:    time.Now().Add(24 * time.Hour).Unix(),
		constants.Type:   "email_confirmation",
	})
	return token.SignedString([]byte(signingKey))
}

func (a *AuthService) SendConfirmationMessage(toEmail string, token string) error {
	link := fmt.Sprintf("http://uchipython.duckdns.org/confirm?token=%s", token)
	messageBody := fmt.Sprintf(`
Click the link to confirm your account:
%s
`, link)

	m := gomail.NewMessage()

	m.SetHeader("From", a.notifyCfg.From)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Confirm your email")
	m.SetBody("text/html", messageBody)

	d := gomail.NewDialer(a.smtpCfg.Host, a.smtpCfg.Port, a.smtpCfg.Email, a.smtpCfg.Password)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("gomail.NewDialer: %w", err)
	}
	return nil
}

func (a *AuthService) Confirm(ctx context.Context, tokenStr string) error {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.hashConfig.SigningKey), nil
	})
	if err != nil {
		return err
	}

	if err != nil || !token.Valid {
		return errlist.ErrUnauthorized
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims[constants.Type] != "email_confirmation" {
		return errlist.ErrUnauthorized
	}

	userID := claims[constants.UserID].(float64)

	return a.repo.UpdateUser(ctx, &model.UserData{
		ID:          int64(userID),
		UpdatedAt:   time.Now(),
		IsConfirmed: true,
	})
}
