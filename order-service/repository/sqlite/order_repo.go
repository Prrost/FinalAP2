package sqlite

import (
	"database/sql"
	"time"

	"github.com/Prrost/FinalAP2/order-service/domain"
	_ "github.com/mattn/go-sqlite3"
)

type orderRepo struct {
	db *sql.DB
}

func NewOrderRepo(db *sql.DB) domain.OrderRepository {
	return &orderRepo{db: db}
}

func (r *orderRepo) Create(o *domain.Order) error {
	stmt := `INSERT INTO orders(user_email, book_id, taken_at, due_at, waiting)
             VALUES (?, ?, ?, ?, ?)`
	res, err := r.db.Exec(stmt, o.UserEmail, o.BookID, o.TakenAt, o.DueAt, o.Waiting)
	if err != nil {
		return err
	}
	o.ID, _ = res.LastInsertId()
	return nil
}

func (r *orderRepo) GetByID(id int64) (*domain.Order, error) {
	row := r.db.QueryRow(
		`SELECT id, user_email, book_id, taken_at, due_at, returned_at, waiting
         FROM orders WHERE id = ?`, id,
	)
	var o domain.Order
	var ret sql.NullTime
	if err := row.Scan(&o.ID, &o.UserEmail, &o.BookID, &o.TakenAt, &o.DueAt, &ret, &o.Waiting); err != nil {
		return nil, err
	}
	if ret.Valid {
		o.ReturnedAt = &ret.Time
	}
	return &o, nil
}

func (r *orderRepo) ListWaiting(bookID int64) ([]*domain.Order, error) {
	rows, err := r.db.Query(
		`SELECT id, user_email, book_id, taken_at, due_at, returned_at, waiting
         FROM orders WHERE book_id = ? AND waiting = 1`, bookID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.Order
	for rows.Next() {
		var o domain.Order
		var ret sql.NullTime
		if err := rows.Scan(&o.ID, &o.UserEmail, &o.BookID, &o.TakenAt, &o.DueAt, &ret, &o.Waiting); err != nil {
			return nil, err
		}
		if ret.Valid {
			o.ReturnedAt = &ret.Time
		}
		list = append(list, &o)
	}
	return list, nil
}

func (r *orderRepo) MarkReturned(id int64, returnedAt time.Time) error {
	_, err := r.db.Exec(
		`UPDATE orders SET returned_at = ?, waiting = 0 WHERE id = ?`,
		returnedAt, id,
	)
	return err
}
