package domain

import "time"

type Order struct {
	ID         int64
	UserID     int64
	BookID     int64
	TakenAt    time.Time
	DueAt      time.Time
	ReturnedAt *time.Time
	Waiting    bool
}

type OrderRepository interface {
	Create(o *Order) error
	GetByID(id int64) (*Order, error)
	ListWaiting(bookID int64) ([]*Order, error)
	MarkReturned(id int64, returnedAt time.Time) error
}
