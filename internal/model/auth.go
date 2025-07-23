package model

import "study_buddy/pkg/constants"

type AuthToken struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}

func NewAuthToken(token string, role int32) *AuthToken {
	var roleStr string
	switch role {
	case constants.Admin:
		roleStr = "admin"
	case constants.User:
		roleStr = "user"
	default:
	}
	return &AuthToken{
		Token: token,
		Role:  roleStr,
	}
}
