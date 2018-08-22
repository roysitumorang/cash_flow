package user

import (
	"time"
)

type (
	User struct {
		Id              int64      `json:"id"`
		Name            string     `json:"name"`
		Email           string     `json:"email"`
		Password        string     `json:"-"`
		PasswordHash    string     `json:"-"`
		PasswordToken   *string    `json:"-"`
		ActivationToken *string    `json:"-"`
		ActivatedAt     *time.Time `json:"activated_at"`
		TimeZone        string     `json:"time_zone"`
		CreatedAt       time.Time  `json:"created_at"`
		UpdatedAt       *time.Time `json:"updated_at"`
		DeletedAt       *time.Time `json:"-"`
	}
	Users []*User
)
