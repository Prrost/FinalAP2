package middleware

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"time"
)

func LoggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()

	res, err := handler(ctx, req)
	if err != nil {
		log.Printf("Completed: %s, Duration: %s, Error: %v\n", info.FullMethod, time.Since(start), err)
	} else {
		log.Printf("Completed: %s, Duration: %s \n", info.FullMethod, time.Since(start))
	}

	return res, err
}
