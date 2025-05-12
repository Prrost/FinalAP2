package grpc

import (
	"github.com/Prrost/FinalAP2/order-service/infra/logger"
	"github.com/Prrost/FinalAP2/order-service/usecase"
	pb "github.com/Prrost/protoFinalAP2/order/order"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

func RunGRPC(uc *usecase.OrderUsecase, port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	srv := grpc.NewServer()

	// Регистрируем свой сервис
	pb.RegisterOrderServiceServer(srv, NewOrderHandler(uc))

	// Регистрируем reflection, чтобы grpcurl и другие клиенты могли узнать API “на лету”
	reflection.Register(srv)

	logger.Log.Infof("gRPC→ listening on %s", port)
	return srv.Serve(lis)
}
