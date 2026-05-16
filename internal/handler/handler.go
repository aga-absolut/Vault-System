package handler

import (
	"context"
	"errors"

	"github.com/aga-absolut/Vault-System/internal/convert"
	"github.com/aga-absolut/Vault-System/internal/errs"
	"github.com/aga-absolut/Vault-System/internal/service/auth"
	"github.com/aga-absolut/Vault-System/internal/service/credentials"
	"github.com/aga-absolut/Vault-System/internal/token"
	pb "github.com/aga-absolut/Vault-System/proto/vault_system"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// VaultSystem implements the gRPC server handlers.
type VaultSystem struct {
	pb.UnimplementedVaultSystemServer
	token              token.Provider
	authService        auth.Service
	credentialsService credentials.Service
}

// NewVaultSystem creates a new VaultSystem handler instance.
func NewVaultSystem(token token.Provider, authService auth.Service, credentialsService credentials.Service) *VaultSystem {
	return &VaultSystem{
		token:              token,
		authService:        authService,
		credentialsService: credentialsService,
	}
}

// Register handles user registration requests.
func (s *VaultSystem) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	token, err := s.authService.RegisterUser(ctx, in.Name, in.Password)
	if err != nil {
		if errors.Is(err, errs.ErrLoginAlreadyUsed) {
			return nil, status.Error(codes.AlreadyExists, "login already used")
		}
		if errors.Is(err, errs.ErrTooShortPassword) {
			return nil, status.Error(codes.InvalidArgument, "too short password")
		}
		return nil, status.Error(codes.Internal, "failed to register")
	}

	return &pb.RegisterResponse{Token: token}, nil
}

// Login handles user authentication requests.
func (s *VaultSystem) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	token, err := s.authService.LoginUser(ctx, in.Name, in.Password)
	if err != nil {
		if errors.Is(err, errs.ErrIncorrectLoginOrPassword) {
			return nil, status.Error(codes.Unauthenticated, "incorrect login or password")
		}
		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &pb.LoginResponse{Token: token}, nil
}

// SetData saves user data on the server.
func (s *VaultSystem) SetData(ctx context.Context, in *pb.SetDataRequest) (*emptypb.Empty, error) {
	userName, err := s.getUserName(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get name from context")
	}

	record := convert.ToRecord(in.Data)
	record.UserName = userName

	err = s.credentialsService.SetData(ctx, record)
	if err != nil {
		if errors.Is(err, errs.ErrMetaAlreadyUsed) {
			return nil, status.Error(codes.Internal, "meta already used:")
		}
		return nil, status.Error(codes.Internal, "failed to save data:")
	}

	return &emptypb.Empty{}, nil
}

// GetData retrieves user data by metadata key.
func (s *VaultSystem) GetData(ctx context.Context, in *pb.GetDataRequest) (*pb.GetDataResponse, error) {
	userName, err := s.getUserName(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get name from context")
	}

	data, err := s.credentialsService.GetData(ctx, userName, in.Meta)
	if err != nil {
		if errors.Is(err, errs.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "record not found")
		}
		return nil, status.Error(codes.Internal, "failed to get data")
	}

	return &pb.GetDataResponse{Data: convert.ToProtoRecord(data)}, nil
}

// GetListMeta returns a list of user record metadata.
func (s *VaultSystem) GetListMeta(ctx context.Context, in *emptypb.Empty) (*pb.GetListMetaResponse, error) {
	userName, err := s.getUserName(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get name from context")
	}

	list, err := s.credentialsService.GetListMeta(ctx, userName)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get list")
	}

	return &pb.GetListMetaResponse{Data: list}, nil
}

// DeleteData removes user data by metadata key.
func (s *VaultSystem) DeleteData(ctx context.Context, in *pb.DeleteDataRequest) (*emptypb.Empty, error) {
	userName, err := s.getUserName(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get name from context")
	}

	if err := s.credentialsService.DeleteData(ctx, userName, in.Meta); err != nil {
		if errors.Is(err, errs.ErrRecordNotFound) {
			return nil, status.Error(codes.Internal, "data is not found")
		}
		return nil, status.Error(codes.Internal, "failed to delete data")
	}
	return &emptypb.Empty{}, nil
}

// UpdateData updates existing user data.
func (s *VaultSystem) UpdateData(ctx context.Context, in *pb.UpdateDataRequest) (*emptypb.Empty, error) {
	userName, err := s.getUserName(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get name from context")
	}

	record := convert.ToRecord(in.Data)
	record.UserName = userName

	if err := s.credentialsService.UpdateData(ctx, record); err != nil {
		return nil, status.Error(codes.Internal, "failed to update data")
	}
	return &emptypb.Empty{}, nil
}

// getUserName extracts the authenticated username from context.
func (s *VaultSystem) getUserName(ctx context.Context) (string, error) {
	username, ok := s.token.UserNameFromContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "user not authenticated")
	}
	return username, nil
}