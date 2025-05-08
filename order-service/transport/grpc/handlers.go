package grpc

import (
	"context"

	"github.com/Prrost/FinalAP2/order-service/infra/logger"
	"github.com/Prrost/FinalAP2/order-service/usecase"
	pb "github.com/Prrost/protoFinalAP2/order/order"
)

type OrderHandler struct {
	pb.UnimplementedOrderServiceServer
	uc *usecase.OrderUsecase
}

func NewOrderHandler(uc *usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{uc: uc}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderReq) (*pb.CreateOrderRes, error) {
	o, err := h.uc.CreateOrder(req.UserId, req.BookId, int(req.DueDays))
	if err != nil {
		logger.Log.Errorf("CreateOrder error: %v", err)
		return nil, err
	}
	return &pb.CreateOrderRes{OrderId: o.ID}, nil
}

func (h *OrderHandler) ReturnOrder(ctx context.Context, req *pb.ReturnOrderReq) (*pb.ReturnOrderRes, error) {
	if err := h.uc.ReturnOrder(req.OrderId); err != nil {
		logger.Log.Errorf("ReturnOrder error: %v", err)
		return nil, err
	}
	return &pb.ReturnOrderRes{}, nil
}
