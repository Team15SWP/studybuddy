package model

type ErrorResponse struct {
	Error string `json:"error"`
}

type AuthInfo struct {
	UserID int64  `db:"user_id"`
	Role   int32  `db:"role"`
	Name   string `db:"username"`
	Email  string `db:"email"`
}
