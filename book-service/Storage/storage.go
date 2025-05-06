package Storage

import "github.com/Prrost/FinalAP2/book-service/domain"

// Storage — интерфейс для работы с БД
type Storage interface {
	ListBooks() ([]domain.Book, error)
	GetBook(id int64) (domain.Book, error)
	CreateBook(b domain.Book) (int64, error)
	UpdateBook(b domain.Book) error
	DeleteBook(id int64) error
}
