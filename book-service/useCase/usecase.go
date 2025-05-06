package useCase

import (
	"github.com/Prrost/FinalAP2/book-service/Storage"
	"github.com/Prrost/FinalAP2/book-service/domain"
)

type UseCase struct {
	storage Storage.Storage
}

func NewUseCase(st Storage.Storage) *UseCase {
	return &UseCase{storage: st}
}

func (uc *UseCase) ListBooks() ([]domain.Book, error) {
	return uc.storage.ListBooks()
}

func (uc *UseCase) GetBook(id int64) (domain.Book, error) {
	return uc.storage.GetBook(id)
}

func (uc *UseCase) CreateBook(b domain.Book) (int64, error) {
	return uc.storage.CreateBook(b)
}

func (uc *UseCase) UpdateBook(b domain.Book) error {
	return uc.storage.UpdateBook(b)
}

func (uc *UseCase) DeleteBook(id int64) error {
	return uc.storage.DeleteBook(id)
}
