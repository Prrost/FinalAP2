package Storage

import (
	"database/sql"
	"log"

	"github.com/Prrost/FinalAP2/book-service/config"
	"github.com/Prrost/FinalAP2/book-service/domain"
	_ "github.com/mattn/go-sqlite3"
)

type sqliteStorage struct {
	db *sql.DB
}

func NewSQLiteStorage(cfg *config.Config) Storage {
	db, err := sql.Open("sqlite3", cfg.DBPath)
	if err != nil {
		log.Fatalf("failed to open sqlite: %v", err)
	}
	// Миграция таблицы
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS books (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT,
        author TEXT,
        isbn TEXT,
        totalQuantity INTEGER,
        availableQuantity INTEGER
    );`)
	if err != nil {
		log.Fatalf("failed to migrate sqlite: %v", err)
	}
	return &sqliteStorage{db: db}
}

func (s *sqliteStorage) ListBooks() ([]domain.Book, error) {
	rows, err := s.db.Query(`SELECT id, title, author, isbn, totalQuantity, availableQuantity FROM books`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Book
	for rows.Next() {
		var b domain.Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.ISBN, &b.TotalQuantity, &b.AvailableQuantity); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, nil
}

func (s *sqliteStorage) GetBook(id int64) (domain.Book, error) {
	var b domain.Book
	err := s.db.QueryRow(
		`SELECT id, title, author, isbn, totalQuantity, availableQuantity FROM books WHERE id = ?`,
		id,
	).Scan(&b.ID, &b.Title, &b.Author, &b.ISBN, &b.TotalQuantity, &b.AvailableQuantity)
	return b, err
}

func (s *sqliteStorage) CreateBook(b domain.Book) (int64, error) {
	res, err := s.db.Exec(
		`INSERT INTO books (title, author, isbn, totalQuantity, availableQuantity) VALUES (?, ?, ?, ?, ?)`,
		b.Title, b.Author, b.ISBN, b.TotalQuantity, b.AvailableQuantity,
	)
	if err != nil {
		log.Printf("failed to create book: %v", err)
		return 0, err
	}
	return res.LastInsertId()
}

func (s *sqliteStorage) UpdateBook(b domain.Book) error {
	_, err := s.db.Exec(
		`UPDATE books SET title=?, author=?, isbn=?, totalQuantity=?, availableQuantity=? WHERE id = ?`,
		b.Title, b.Author, b.ISBN, b.TotalQuantity, b.AvailableQuantity, b.ID,
	)
	return err
}

func (s *sqliteStorage) DeleteBook(id int64) error {
	_, err := s.db.Exec(`DELETE FROM books WHERE id = ?`, id)
	return err
}
