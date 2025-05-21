package useCase

import (
	"errors"
	"testing"

	"github.com/Prrost/FinalAP2/book-service/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStorage реализует domain.Storage и используется для unit-тестов
type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) ListBooks() ([]domain.Book, error) {
	args := m.Called()
	return args.Get(0).([]domain.Book), args.Error(1)
}

func (m *MockStorage) GetBook(id int64) (domain.Book, error) {
	args := m.Called(id)
	return args.Get(0).(domain.Book), args.Error(1)
}

func (m *MockStorage) CreateBook(b domain.Book) (int64, error) {
	args := m.Called(b)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockStorage) UpdateBook(b domain.Book) error {
	args := m.Called(b)
	return args.Error(0)
}

func (m *MockStorage) DeleteBook(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

// --------------------------
// Позитивные кейсы
// --------------------------

func TestUseCase_ListBooks(t *testing.T) {
	mockSt := new(MockStorage)
	expected := []domain.Book{
		{ID: 1, Title: "Test Book"},
	}
	// важно: возвращаем именно этот слайс и nil
	mockSt.On("ListBooks").Return(expected, nil)

	uc := NewUseCase(mockSt)
	result, err := uc.ListBooks()

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockSt.AssertExpectations(t)
}

func TestUseCase_GetBook(t *testing.T) {
	mockSt := new(MockStorage)
	sample := domain.Book{ID: 2, Title: "Another Book"}
	mockSt.On("GetBook", int64(2)).Return(sample, nil)

	uc := NewUseCase(mockSt)
	result, err := uc.GetBook(2)

	assert.NoError(t, err)
	assert.Equal(t, sample, result)
	mockSt.AssertExpectations(t)
}

func TestUseCase_CreateBook(t *testing.T) {
	mockSt := new(MockStorage)
	newBook := domain.Book{Title: "New Book"}
	mockSt.On("CreateBook", newBook).Return(int64(100), nil)

	uc := NewUseCase(mockSt)
	id, err := uc.CreateBook(newBook)

	assert.NoError(t, err)
	assert.Equal(t, int64(100), id)
	mockSt.AssertExpectations(t)
}

func TestUseCase_UpdateBook(t *testing.T) {
	mockSt := new(MockStorage)
	updBook := domain.Book{ID: 5, Title: "Updated Title"}
	mockSt.On("UpdateBook", updBook).Return(nil)

	uc := NewUseCase(mockSt)
	err := uc.UpdateBook(updBook)

	assert.NoError(t, err)
	mockSt.AssertExpectations(t)
}

func TestUseCase_DeleteBook(t *testing.T) {
	mockSt := new(MockStorage)
	mockSt.On("DeleteBook", int64(5)).Return(nil)

	uc := NewUseCase(mockSt)
	err := uc.DeleteBook(5)

	assert.NoError(t, err)
	mockSt.AssertExpectations(t)
}

// --------------------------
// Негативные кейсы
// --------------------------

func TestUseCase_ListBooks_Error(t *testing.T) {
	mockSt := new(MockStorage)
	// возвращаем пустой слайс (не nil!) и ошибку
	mockSt.On("ListBooks").Return([]domain.Book{}, errors.New("db error"))

	uc := NewUseCase(mockSt)
	result, err := uc.ListBooks()

	assert.Error(t, err)
	assert.Empty(t, result) // проверяем, что слайс пуст
	mockSt.AssertExpectations(t)
}

func TestUseCase_GetBook_Error(t *testing.T) {
	mockSt := new(MockStorage)
	mockSt.On("GetBook", int64(1)).Return(domain.Book{}, errors.New("not found"))

	uc := NewUseCase(mockSt)
	result, err := uc.GetBook(1)

	assert.Error(t, err)
	assert.Empty(t, result)
	mockSt.AssertExpectations(t)
}

func TestUseCase_CreateBook_Error(t *testing.T) {
	mockSt := new(MockStorage)
	bk := domain.Book{Title: "ErrBook"}
	mockSt.On("CreateBook", bk).Return(int64(0), errors.New("cannot create"))

	uc := NewUseCase(mockSt)
	id, err := uc.CreateBook(bk)

	assert.Error(t, err)
	assert.Zero(t, id)
	mockSt.AssertExpectations(t)
}

func TestUseCase_UpdateBook_Error(t *testing.T) {
	mockSt := new(MockStorage)
	bk := domain.Book{ID: 7, Title: "ErrUpdate"}
	mockSt.On("UpdateBook", bk).Return(errors.New("cannot update"))

	uc := NewUseCase(mockSt)
	err := uc.UpdateBook(bk)

	assert.Error(t, err)
	mockSt.AssertExpectations(t)
}

func TestUseCase_DeleteBook_Error(t *testing.T) {
	mockSt := new(MockStorage)
	mockSt.On("DeleteBook", int64(7)).Return(errors.New("cannot delete"))

	uc := NewUseCase(mockSt)
	err := uc.DeleteBook(7)

	assert.Error(t, err)
	mockSt.AssertExpectations(t)
}
