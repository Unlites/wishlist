package domain

import "time"

type Wish struct {
	Id          int
	UserId      int
	Title       string
	Description string
	IsReserved  *bool
	ReservedBy  *int
	CreatedAt   time.Time
}

func (w *Wish) SetReserved(isReserved bool, reservedBy int) {
	w.IsReserved = &isReserved

	if !isReserved {
		w.ReservedBy = nil
	} else {
		w.ReservedBy = &reservedBy
	}
}
