package model

import "time"

type Notification struct {
	UserID  int64     `json:"user_id"`
	Enabled bool      `json:"enabled"`
	Time24  time.Time `json:"time_24"`
	Days    []int     `json:"days"`
}
