package auth

import (
	"context"

	"github.com/aga-absolut/Vault-System/internal/token"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UnaryServerInterceptor(tokenProvider token.Provider) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if isPublicMethod(info.FullMethod) {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
		}

		authHeader := md.Get("authorization")
		if len(authHeader) == 0 {
			return nil, status.Error(codes.Unauthenticated, "authorization token is required")
		}

		t := authHeader[0]
		userName, err := tokenProvider.ValidateToken(t)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid or expired token")
		}

		newCtx := context.WithValue(ctx, token.UserIDKey{}, userName)

		return handler(newCtx, req)
	}
}

func isPublicMethod(fullMethod string) bool {
	public := []string{
		"/vault.VaultSystem/Register",
		"/vault.VaultSystem/Login",
	}

	for _, m := range public {
		if fullMethod == m {
			return true
		}
	}
	return false
}
