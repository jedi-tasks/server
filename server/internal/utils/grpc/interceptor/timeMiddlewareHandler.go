package interceptor

import (
	"context"
	"server/internal/metrics"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

type TimeInterceptor struct{}

func NewTimeInterceptor() *TimeInterceptor {
	return &TimeInterceptor{}
}

func (i *TimeInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {

		// Фиксируем время старта
		start := time.Now()

		// Вызываем следующий хэндлер
		res, err := handler(ctx, req)

		// Вычисляем продолжительность
		duration := time.Since(start)

		isSuccess := err == nil

		// Записываем информацию о времени ответа с использованием прометеуса
		metrics.GlobalMetrics.ResponseTimeMetric.WithLabelValues(
			info.FullMethod,
			strconv.FormatBool(isSuccess),
		).Observe(duration.Seconds())

		return res, err
	}
}

func (i *TimeInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

		// Фиксируем время старта
		start := time.Now()

		// Вызываем следующий хэндлер
		err := handler(srv, stream)

		// Вычисляем продолжительность
		duration := time.Since(start)

		isSuccess := err == nil

		// Записываем информацию о времени ответа с использованием прометеуса
		metrics.GlobalMetrics.ResponseTimeMetric.WithLabelValues(
			info.FullMethod,
			strconv.FormatBool(isSuccess),
		).Observe(duration.Seconds())

		return err
	}
}
