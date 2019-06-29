package tests

import (
	"time"
)

type UserData struct {
	UserId      string
	HashedToken string
}

type Shift struct {
	ID        string
	StartDate time.Time
	EndDate   time.Time
	Slots     int
	Hireds    []string
	Deadline  bool
}
