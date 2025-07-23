package middlewares

import (
	"fmt"

	"study_buddy/internal/config"
	"study_buddy/internal/model"
	"study_buddy/pkg/constants"
	"study_buddy/pkg/errlist"

	"github.com/dgrijalva/jwt-go"
)

func ParseToken(tokenString string, hashCfg *config.HashConfig) (*model.AuthInfo, error) {
	response := model.AuthInfo{}

	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v, %w", token.Header["alg"], errlist.ErrUnauthorized)
		}
		return []byte(hashCfg.SigningKey), nil
	})
	if err != nil {
		return &response, errlist.ErrUnauthorized
	}

	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok {
		return &response, fmt.Errorf("getting custom claims from token 1: %w", errlist.ErrUnauthorized)
	}

	var userId, role float64
	if userId, ok = (*claims)[constants.UserID].(float64); !ok {
		return &response, fmt.Errorf("getting custom claims from token 2: %w", errlist.ErrUnauthorized)
	}
	response.UserID = int64(userId)

	if role, ok = (*claims)[constants.Role].(float64); !ok {
		return &response, fmt.Errorf("getting custom claims from token 3: %w", errlist.ErrUnauthorized)
	}
	response.Role = int32(role)

	if response.Name, ok = (*claims)[constants.Name].(string); !ok {
		return &response, fmt.Errorf("getting custom claims from token 4: %w", errlist.ErrUnauthorized)
	}

	if response.Email, ok = (*claims)[constants.Email].(string); !ok {
		return &response, fmt.Errorf("getting custom claims from token 5: %w", errlist.ErrUnauthorized)
	}

	return &response, nil
}
