package usecase

import (
	"fmt"
	"log"
	"time"

	"github.com/Prrost/FinalAP2/order-service/domain"
	"github.com/Prrost/FinalAP2/order-service/infra/email"
	"github.com/Prrost/FinalAP2/order-service/publisher"
)

type OrderUsecase struct {
	repo        domain.OrderRepository
	RMQ         *publisher.OrderCreatedPublisher
	EmailSender email.Sender
}

// Теперь конструктор принимает email.Sender и возвращает ошибку, если RMQ не поднялся
func NewOrderUC(r domain.OrderRepository, emailSender email.Sender) (*OrderUsecase, error) {
	rmq, err := publisher.NewOrderCreatedPublisher("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, fmt.Errorf("failed to init RMQ publisher: %w", err)
	}
	return &OrderUsecase{
		repo:        r,
		RMQ:         rmq,
		EmailSender: emailSender,
	}, nil
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

	// 📨 Публикуем в RabbitMQ
	if err := u.RMQ.OrderCreatedPublish(userEmail); err != nil {
		return nil, err
	}

	// ✉️ Отправляем email пользователю
	subject := "Order Confirmation"
	body := fmt.Sprintf(
		"Hello!\n\nYour order for book #%d has been created.\nDue date: %s",
		o.BookID, o.DueAt.Format("02 Jan 2006"),
	)
	if err := u.EmailSender.Send(userEmail, subject, body); err != nil {
		log.Printf("[Email] Failed to send: %v", err)
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

func (u *OrderUsecase) GetByID(id int64) (*domain.Order, error) {
	return u.repo.GetByID(id)
}

func (u *OrderUsecase) ListWaiting(bookID int64) ([]*domain.Order, error) {
	return u.repo.ListWaiting(bookID)
}
