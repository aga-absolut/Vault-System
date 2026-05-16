package auth

import (
	"context"
	"errors"

	"github.com/aga-absolut/Vault-System/internal/errs"
	"github.com/aga-absolut/Vault-System/internal/interceptors/logger"
	"github.com/aga-absolut/Vault-System/internal/storage/postgres"
	"github.com/aga-absolut/Vault-System/internal/token"
	"golang.org/x/crypto/bcrypt"
)

// Service defines authentication service methods.
type Service interface {
	RegisterUser(context.Context, string, string) (string, error)
	LoginUser(context.Context, string, string) (string, error)
}

// service implements authentication business logic.
type service struct {
	log   *logger.Logger
	db    postgres.Storage
	token token.Provider
}

// NewService creates a new authentication service instance.
func NewService(log *logger.Logger, db postgres.Storage, token token.Provider) Service {
	return &service{
		log:   log,
		db:    db,
		token: token,
	}
}

// RegisterUser creates a new user and returns a JWT token.
func (s *service) RegisterUser(ctx context.Context, name, password string) (string, error) {
	if name == "" || password == "" {
		return "", errs.ErrIncorrectLoginOrPassword
	}

	if len(password) < 8 {
		return "", errs.ErrTooShortPassword
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		s.log.Errorw("failed to hash password", "error", err)
		return "", errs.ErrInternal
	}

	if err := s.db.AddUser(ctx, name, string(hashPassword)); err != nil {
		s.log.Errorw("failed to add user", "username", name, "error", err)
		if errors.Is(err, errs.ErrLoginAlreadyUsed) {
			return "", err
		}
		return "", errs.ErrInternal
	}

	tokenStr, err := s.token.BuildJWTString(name)
	if err != nil {
		s.log.Errorw("failed to generate token", "username", name, "error", err)
		return "", errs.ErrInternal
	}

	return tokenStr, nil
}

// LoginUser authenticates a user and returns a JWT token.
func (s *service) LoginUser(ctx context.Context, name, password string) (string, error) {
	hashPassword, err := s.db.CheckUser(ctx, name)
	if err != nil {
		if errors.Is(err, errs.ErrIncorrectLoginOrPassword) {
			s.log.Errorw("login attempt failed", "username", name, "error", err)
			return "", errs.ErrIncorrectLoginOrPassword
		}

		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password)); err != nil {
		s.log.Errorw("invalid password", "username", name)
		return "", errs.ErrIncorrectLoginOrPassword
	}

	tokenStr, err := s.token.BuildJWTString(name)
	if err != nil {
		s.log.Errorw("failed to generate token", "username", name, "error", err)
		return "", errs.ErrInternal
	}

	return tokenStr, nil
}
