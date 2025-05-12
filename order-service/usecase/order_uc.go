package usecase

import (
	"github.com/Prrost/FinalAP2/order-service/publisher"
	"log"
	"time"

	"github.com/Prrost/FinalAP2/order-service/domain"
)

type OrderUsecase struct {
	repo domain.OrderRepository
	RMQ  *publisher.OrderCreatedPublisher
}

func NewOrderUC(r domain.OrderRepository) *OrderUsecase {
	OCPublisher, err := publisher.NewOrderCreatedPublisher("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
		return &OrderUsecase{}
	}
	return &OrderUsecase{repo: r, RMQ: OCPublisher}
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
	err := u.RMQ.OrderCreatedPublish(userEmail)
	if err != nil {
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
