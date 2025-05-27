package handlers

import (
	"context"
	"log"

	bookpb "github.com/Prrost/protoFinalAP2/books"
	emptypb "google.golang.org/protobuf/types/known/emptypb"

	"github.com/Prrost/FinalAP2/book-service/domain"
	"github.com/Prrost/FinalAP2/book-service/useCase"
)

type Server struct {
	bookpb.UnimplementedBookServiceServer
	uc *useCase.UseCase
}

func NewServer(uc *useCase.UseCase) *Server {
	return &Server{uc: uc}
}

func (s *Server) ListBooks(ctx context.Context, _ *emptypb.Empty) (*bookpb.BookList, error) {
	log.Println("Entery")
	books, err := s.uc.ListBooks()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	log.Println("Ok")
	resp := &bookpb.BookList{}
	for _, b := range books {
		resp.Books = append(resp.Books, &bookpb.Book{
			Id:                b.ID,
			Title:             b.Title,
			Author:            b.Author,
			Isbn:              b.ISBN,
			TotalQuantity:     b.TotalQuantity,
			AvailableQuantity: b.AvailableQuantity,
		})
	}
	return resp, nil
}

func (s *Server) GetBook(ctx context.Context, req *bookpb.BookId) (*bookpb.Book, error) {
	b, err := s.uc.GetBook(req.Id)
	if err != nil {
		return nil, err
	}
	return &bookpb.Book{
		Id:                b.ID,
		Title:             b.Title,
		Author:            b.Author,
		Isbn:              b.ISBN,
		TotalQuantity:     b.TotalQuantity,
		AvailableQuantity: b.AvailableQuantity,
	}, nil
}

func (s *Server) CreateBook(ctx context.Context, req *bookpb.Book) (*bookpb.Book, error) {
	dom := domain.Book{
		Title:             req.Title,
		Author:            req.Author,
		ISBN:              req.Isbn,
		TotalQuantity:     req.TotalQuantity,
		AvailableQuantity: req.AvailableQuantity,
	}
	id, err := s.uc.CreateBook(dom)
	if err != nil {
		return nil, err
	}
	req.Id = id
	return req, nil
}

func (s *Server) UpdateBook(ctx context.Context, req *bookpb.Book) (*bookpb.Book, error) {
	dom := domain.Book{
		ID:                req.Id,
		Title:             req.Title,
		Author:            req.Author,
		ISBN:              req.Isbn,
		TotalQuantity:     req.TotalQuantity,
		AvailableQuantity: req.AvailableQuantity,
	}
	if err := s.uc.UpdateBook(dom); err != nil {
		return nil, err
	}
	return req, nil
}

func (s *Server) DeleteBook(ctx context.Context, req *bookpb.BookId) (*emptypb.Empty, error) {
	if err := s.uc.DeleteBook(req.Id); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
