package logger

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// Logger wraps a zap sugared logger instance.
type Logger struct {
	*zap.SugaredLogger
}

// NewLogger creates a new development logger instance.
func NewLogger() *Logger {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		panic("failed to create logger: " + err.Error())
	}

	return &Logger{
		SugaredLogger: zapLogger.Sugar(),
	}
}

// UnaryServerInterceptor logs information about incoming gRPC requests.
func (l *Logger) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(start)
		st, _ := status.FromError(err)

		l.Infoln(
			"\n",
			"-----GRPC REQUEST-----\n",
			"Method:", info.FullMethod, "\n",
			"Status:", st.Message(), "\n",
			"Duration:", duration, "\n",
			"Status code:", st.Code(), "\n",
		)

		return resp, err
	}
}