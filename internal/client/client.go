package client

import (
	"context"
	"fmt"

	"github.com/aga-absolut/Vault-System/internal/convert"
	"github.com/aga-absolut/Vault-System/internal/errs"
	"github.com/aga-absolut/Vault-System/internal/models"
	pb "github.com/aga-absolut/Vault-System/proto/vault_system"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Client provides methods for interacting with the VaultSystem gRPC server.
type Client struct {
	client pb.VaultSystemClient
	conn   *grpc.ClientConn

	// Token stores the authentication token for authorized requests.
	Token string
}

// NewClient creates a new gRPC client connection to the server.
func NewClient(serverAddr string) *Client {
	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(
		insecure.NewCredentials(),
	))
	if err != nil {
		panic(err)
	}

	client := pb.NewVaultSystemClient(conn)
	return &Client{
		client: client,
		conn:   conn,
	}
}

// Register creates a new user account and stores the received auth token.
func (c *Client) Register(ctx context.Context, userName, password string) error {
	resp, err := c.client.Register(ctx, &pb.RegisterRequest{Name: userName, Password: password})
	if err != nil {
		if status, ok := status.FromError(err); ok {
			switch status.Code() {
			case codes.AlreadyExists:
				return errs.ErrLoginAlreadyUsed
			case codes.InvalidArgument:
				return errs.ErrTooShortPassword
			case codes.Internal:
				return errs.ErrInternal
			case codes.Unavailable:
				return fmt.Errorf("сервер недоступен (не запущен?)")
			default:
				return fmt.Errorf("gRPC ошибка: %s (%s)", status.Code(), status.Message())
			}
		}
		return fmt.Errorf("register failed: %w", err)
	}

	c.Token = resp.Token
	return nil
}

// Login authenticates the user and stores the received auth token.
func (c *Client) Login(ctx context.Context, userName, password string) error {
	resp, err := c.client.Login(ctx, &pb.LoginRequest{Name: userName, Password: password})
	if err != nil {
		if status, ok := status.FromError(err); ok {
			switch status.Code() {
			case codes.Unauthenticated:
				return errs.ErrIncorrectLoginOrPassword
			case codes.Internal:
				return errs.ErrInternal
			}
		}
		return err
	}

	c.Token = resp.Token
	return nil
}

// SetData sends a new record to the server for storage.
func (c *Client) SetData(ctx context.Context, recordType, recordMeta string, recordData []byte) error {
	record := &models.Record{
		Type: recordType,
		Meta: recordMeta,
		Data: recordData,
	}

	ctx = withToken(ctx, c.Token)
	_, err := c.client.SetData(ctx, &pb.SetDataRequest{Data: convert.ToProtoRecord(record)})
	if err != nil {
		if status, ok := status.FromError(err); ok {
			switch status.Code() {
			case codes.Internal:
				return errs.ErrInternal
			}
		}
		return err
	}

	return nil
}

// GetData retrieves a record from the server by metadata key.
func (c *Client) GetData(ctx context.Context, meta string) (*models.Record, error) {
	ctx = withToken(ctx, c.Token)
	resp, err := c.client.GetData(ctx, &pb.GetDataRequest{Meta: meta})
	if err != nil {
		if status, ok := status.FromError(err); ok {
			switch status.Code() {
			case codes.Internal:
				return nil, errs.ErrInternal
			case codes.NotFound:
				return nil, errs.ErrRecordNotFound
			}
		}
		return nil, err
	}

	record := convert.ToRecord(resp.Data)
	return record, nil
}

// GetListMeta returns a list of all stored record metadata keys.
func (c *Client) GetListMeta(ctx context.Context) ([]string, error) {
	ctx = withToken(ctx, c.Token)
	resp, err := c.client.GetListMeta(ctx, &emptypb.Empty{})
	if err != nil {
		if status, ok := status.FromError(err); ok {
			switch status.Code() {
			case codes.Internal:
				return nil, errs.ErrInternal
			}
		}
		return nil, err
	}
	return resp.Data, nil
}

// UpdateData updates an existing record on the server.
func (c *Client) UpdateData(ctx context.Context, recordType, recordMeta string, recordData []byte) error {
	record := &models.Record{
		Type: recordType,
		Meta: recordMeta,
		Data: recordData,
	}

	ctx = withToken(ctx, c.Token)
	_, err := c.client.UpdateData(ctx, &pb.UpdateDataRequest{Data: convert.ToProtoRecord(record)})
	if err != nil {
		if status, ok := status.FromError(err); ok {
			switch status.Code() {
			case codes.Internal:
				return errs.ErrInternal
			}
		}
		return err
	}
	return nil
}

// DeleteData removes a record from the server by metadata key.
func (c *Client) DeleteData(ctx context.Context, meta string) error {
	ctx = withToken(ctx, c.Token)
	_, err := c.client.DeleteData(ctx, &pb.DeleteDataRequest{Meta: meta})
	if err != nil {
		if status, ok := status.FromError(err); ok {
			switch status.Code() {
			case codes.Internal:
				return errs.ErrInternal
			}
		}
		return err
	}

	return nil
}

// Close closes the gRPC client connection.
func (c *Client) Close() error {
	return c.conn.Close()
}

// withToken attaches the authentication token to the outgoing request context.
func withToken(ctx context.Context, Token string) context.Context {
	if Token == "" {
		return ctx
	}
	md := metadata.Pairs("authorization", Token)
	return metadata.NewOutgoingContext(ctx, md)
}
