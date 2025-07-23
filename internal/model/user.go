package model

import "time"

type UserData struct {
	ID          int64     `db:"id"`
	Role        int32     `db:"role"`
	Name        string    `db:"username"`
	Email       string    `db:"email"`
	Password    string    `db:"password"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	IsConfirmed bool      `db:"is_confirmed"`
}

type User struct {
	Name  string
	Email string
}
