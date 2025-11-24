package domain

import "time"

type User struct {
	Id       uint64
	Username string
	Email    string
	Password string
	CTime    time.Time
}
