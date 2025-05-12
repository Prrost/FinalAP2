package usecase

import (
	"time"

	"github.com/Prrost/FinalAP2/order-service/domain"
)

type OrderUsecase struct {
	repo domain.OrderRepository
}

func NewOrderUC(r domain.OrderRepository) *OrderUsecase {
	return &OrderUsecase{repo: r}
}

func (u *OrderUsecase) CreateOrder(userEmail string, bookID int64, dueDays int) (*domain.Order, error) {
	o := &domain.Order{
		UserEmail: userEmail,
		BookID:    bookID,
		TakenAt:   time.Now(),
		DueAt:     time.Now().Add(time.Duration(dueDays) * 24 * time.Hour),
		Waiting:   false,
	}
	if err := u.repo.Create(o); err != nil {
		return nil, err
	}
	return o, nil
}

func (u *OrderUsecase) ReturnOrder(id int64) (*domain.Order, error) {
	o, err := u.repo.GetByID(id)
	if err != nil {
		return &domain.Order{}, err
	}
	now := time.Now()
	if err := u.repo.MarkReturned(id, now); err != nil {
		return &domain.Order{}, err
	}
	return o, nil
}
