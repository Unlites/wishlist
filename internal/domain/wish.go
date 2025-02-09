package domain

import "time"

type Wish struct {
	Id          int
	UserId      int
	Title       string
	Description string
	IsReserved  *bool
	CreatedAt   time.Time
}
