package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/aga-absolut/Vault-System/internal/config"
	"github.com/aga-absolut/Vault-System/internal/crypto"
	"github.com/aga-absolut/Vault-System/internal/handler"
	"github.com/aga-absolut/Vault-System/internal/interceptors/auth"
	"github.com/aga-absolut/Vault-System/internal/interceptors/logger"
	"github.com/aga-absolut/Vault-System/internal/interceptors/recovery"
	authService "github.com/aga-absolut/Vault-System/internal/service/auth"
	"github.com/aga-absolut/Vault-System/internal/service/credentials"
	migrate "github.com/aga-absolut/Vault-System/internal/storage"
	"github.com/aga-absolut/Vault-System/internal/storage/postgres"
	"github.com/aga-absolut/Vault-System/internal/token"
	pb "github.com/aga-absolut/Vault-System/proto/vault_system"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.MustLoadConfig()
	log := logger.NewLogger()
	cipher := crypto.NewCipher(cfg.EncryptionKey)
	tokenProvider := token.NewJWTProvider(cfg)

	storage := postgres.NewStorage(cfg)
	// migrate.MustDownMigrations(cfg, log)
	migrate.MustInitMigrations(cfg, log)

	credService := credentials.NewService(log, storage, cipher)
	authService := authService.NewService(log, storage, tokenProvider)

	handler := handler.NewVaultSystem(tokenProvider, authService, credService)

	server := grpc.NewServer(grpc.ChainUnaryInterceptor(
		log.UnaryServerInterceptor(),
		auth.UnaryServerInterceptor(tokenProvider),
		recovery.UnaryServerInterceptor(),
	))

	pb.RegisterVaultSystemServer(server, handler)
	reflection.Register(server)

	listen, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		log.Fatalw("failed to listen on gRPC port", "port", cfg.Address, "error", err)
	}

	go func() {
		log.Infow("starting gRPC server", "addr", cfg.Address)

		if err := server.Serve(listen); err != nil {
			log.Errorw("gRPC server stopped", "error", err)
		}
	}()

	<-ctx.Done()
	log.Info("signal recieved")

	server.GracefulStop()
	log.Info("application stopped successfully")
}
