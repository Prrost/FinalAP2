// integration_test.go
package main

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/emptypb"

	sqlite "github.com/Prrost/FinalAP2/book-service/Storage"
	"github.com/Prrost/FinalAP2/book-service/config"
	"github.com/Prrost/FinalAP2/book-service/internal/handlers"
	"github.com/Prrost/FinalAP2/book-service/useCase"
	bookpb "github.com/Prrost/protoFinalAP2/books"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	// Поднимаем in-memory gRPC-сервер с bufconn
	lis = bufconn.Listen(bufSize)

	// Конфигурация SQLite в памяти
	cfg := &config.Config{
		DBPath: ":memory:",
	}
	storage := sqlite.NewSQLiteStorage(cfg)
	uc := useCase.NewUseCase(storage)

	// Регистрируем gRPC-сервер
	srv := grpc.NewServer()
	bookpb.RegisterBookServiceServer(srv, handlers.NewServer(uc))

	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("bufconn serve error: %v", err)
		}
	}()
}

// bufDialer создаёт соединение к in-memory серверу
func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

// TestCreateAndGetBookGRPC — базовый сценарий: CreateBook → GetBook
func TestCreateAndGetBookGRPC(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(time.Second),
	)
	require.NoError(t, err)
	defer conn.Close()

	client := bookpb.NewBookServiceClient(conn)

	// 1) создаём книгу
	createReq := &bookpb.Book{
		Title:             "Integration Book",
		Author:            "Tester",
		Isbn:              "123-456",
		TotalQuantity:     10,
		AvailableQuantity: 10,
	}
	createResp, err := client.CreateBook(ctx, createReq)
	require.NoError(t, err, "CreateBook должен вернуть без ошибки")
	require.NotZero(t, createResp.Id, "ID должен быть ненулевой")

	// 2) получаем её обратно
	getResp, err := client.GetBook(ctx, &bookpb.BookId{Id: createResp.Id})
	require.NoError(t, err, "GetBook должен вернуть без ошибки")
	require.Equal(t, createReq.Title, getResp.Title, "Title должен совпадать")
	require.Equal(t, createReq.Author, getResp.Author, "Author должен совпадать")
	require.Equal(t, createReq.Isbn, getResp.Isbn, "ISBN должен совпадать")
}

// TestListBooksGRPC — создаём несколько книг и проверяем ListBooks
func TestListBooksGRPC(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(time.Second),
	)
	require.NoError(t, err)
	defer conn.Close()

	client := bookpb.NewBookServiceClient(conn)

	// создаём две книги
	_, err = client.CreateBook(ctx, &bookpb.Book{Title: "Book #1"})
	require.NoError(t, err)
	_, err = client.CreateBook(ctx, &bookpb.Book{Title: "Book #2"})
	require.NoError(t, err)

	// вызываем ListBooks
	listResp, err := client.ListBooks(ctx, &emptypb.Empty{})
	require.NoError(t, err, "ListBooks должен вернуть без ошибки")
	require.GreaterOrEqual(t, len(listResp.Books), 2, "должно быть минимум 2 книги")
}

// TestUpdateBookGRPC — проверяем, что UpdateBook обновляет данные
func TestUpdateBookGRPC(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(time.Second),
	)
	require.NoError(t, err)
	defer conn.Close()

	client := bookpb.NewBookServiceClient(conn)

	// создаём книгу
	createResp, err := client.CreateBook(ctx, &bookpb.Book{Title: "Old Title"})
	require.NoError(t, err)

	// обновляем её
	_, err = client.UpdateBook(ctx, &bookpb.Book{
		Id:    createResp.Id,
		Title: "New Title",
	})
	require.NoError(t, err, "UpdateBook должен вернуть без ошибки")

	// проверяем, что изменилось
	getResp, err := client.GetBook(ctx, &bookpb.BookId{Id: createResp.Id})
	require.NoError(t, err)
	require.Equal(t, "New Title", getResp.Title, "Title должен совпадать")
}

// TestDeleteBookGRPC — проверяем, что после удаления GetBook возвращает ошибку
func TestDeleteBookGRPC(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(time.Second),
	)
	require.NoError(t, err)
	defer conn.Close()

	client := bookpb.NewBookServiceClient(conn)

	// создаём книгу
	createResp, err := client.CreateBook(ctx, &bookpb.Book{Title: "To be deleted"})
	require.NoError(t, err)

	// удаляем
	_, err = client.DeleteBook(ctx, &bookpb.BookId{Id: createResp.Id})
	require.NoError(t, err, "DeleteBook должен вернуть без ошибки")

	// пытаемся получить — ждём ошибку
	_, err = client.GetBook(ctx, &bookpb.BookId{Id: createResp.Id})
	require.Error(t, err, "GetBook после DeleteBook должен вернуть ошибку")
}

// TestGetBook_NotFoundGRPC — негативный кейс: несуществующий ID
func TestGetBook_NotFoundGRPC(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(time.Second),
	)
	require.NoError(t, err)
	defer conn.Close()

	client := bookpb.NewBookServiceClient(conn)

	// пробуем получить несуществующий ID
	_, err = client.GetBook(ctx, &bookpb.BookId{Id: 999999})
	require.Error(t, err, "ожидаем ошибку при запросе несуществующего ID")
}
