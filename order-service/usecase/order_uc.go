package usecase

import (
	"time"

	"github.com/Prrost/FinalAP2/order-service/domain"
	"github.com/Prrost/FinalAP2/order-service/infra/logger"
	natspub "github.com/Prrost/FinalAP2/order-service/infra/nats"
)

type OrderUsecase struct {
	repo domain.OrderRepository
	pub  *natspub.Publisher
}

func NewOrderUC(r domain.OrderRepository, p *natspub.Publisher) *OrderUsecase {
	return &OrderUsecase{repo: r, pub: p}
}

func (u *OrderUsecase) CreateOrder(userID, bookID int64, dueDays int) (*domain.Order, error) {
	o := &domain.Order{
		UserID:  userID,
		BookID:  bookID,
		TakenAt: time.Now(),
		DueAt:   time.Now().Add(time.Duration(dueDays) * 24 * time.Hour),
		Waiting: false,
	}
	if err := u.repo.Create(o); err != nil {
		return nil, err
	}
	evt := struct {
		OrderID int64 `json:"order_id"`
		BookID  int64 `json:"book_id"`
	}{o.ID, o.BookID}
	if err := u.pub.Publish("order.created", evt); err != nil {
		logger.Log.Errorf("failed to publish order.created: %v", err)
	}
	return o, nil
}

func (u *OrderUsecase) ReturnOrder(id int64) error {
	o, err := u.repo.GetByID(id)
	if err != nil {
		return err
	}
	now := time.Now()
	if err := u.repo.MarkReturned(id, now); err != nil {
		return err
	}
	evt := struct {
		OrderID int64 `json:"order_id"`
		BookID  int64 `json:"book_id"`
	}{o.ID, o.BookID}
	if err := u.pub.Publish("order.returned", evt); err != nil {
		logger.Log.Errorf("failed to publish order.returned: %v", err)
	}
	return nil
}
