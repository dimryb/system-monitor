package interceptors

import (
	"context"
	"time"

	i "github.com/dimryb/system-monitor/internal/interface"
	"google.golang.org/grpc"
)

func UnaryLoggerInterceptor(log i.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		log.Infof("gRPC request: %s", info.FullMethod)
		log.Debugf("Request payload: %+v", req)

		resp, err := handler(ctx, req)

		duration := time.Since(start).Milliseconds()
		status := "success"
		if err != nil {
			status = "error"
		}

		log.Infof("gRPC finished: method=%s duration=%dms status=%s error=%v",
			info.FullMethod, duration, status, err)

		return resp, err
	}
}
