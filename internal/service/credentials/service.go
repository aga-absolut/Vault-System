package credentials

import (
	"context"
	"errors"

	"github.com/aga-absolut/Vault-System/internal/crypto"
	"github.com/aga-absolut/Vault-System/internal/errs"
	"github.com/aga-absolut/Vault-System/internal/interceptors/logger"
	"github.com/aga-absolut/Vault-System/internal/models"
	"github.com/aga-absolut/Vault-System/internal/storage/postgres"
)

// Service defines methods for managing encrypted user data.
type Service interface {
	SetData(context.Context, *models.Record) error
	GetData(context.Context, string, string) (*models.Record, error)
	UpdateData(ctx context.Context, record *models.Record) error
	DeleteData(context.Context, string, string) error
	GetListMeta(context.Context, string) ([]string, error)
}

// service implements credential management business logic.
type service struct {
	log    *logger.Logger
	db     postgres.Storage
	cipher *crypto.Cipher
}

// NewService creates a new credentials service instance.
func NewService(log *logger.Logger, db postgres.Storage, cipher *crypto.Cipher) Service {
	return &service{
		log:    log,
		db:     db,
		cipher: cipher,
	}
}

// SetData encrypts and stores user data.
func (s *service) SetData(ctx context.Context, record *models.Record) error {
	encrypted, err := s.cipher.Encrypt(record.Data)
	if err != nil {
		s.log.Errorw("failed to encrypt data", "user", record.UserName, "error", err)
		return errs.ErrInternal
	}

	err = s.db.SetData(ctx, record.Type, record.Meta, record.UserName, encrypted)
	if err != nil {
		if errors.Is(err, errs.ErrMetaAlreadyUsed) {
			return err
		}
		s.log.Errorw("failed to save credential", "user", record.UserName, "error", err)
		return err
	}

	return nil
}

// GetData retrieves and decrypts user data.
func (s *service) GetData(ctx context.Context, userName, meta string) (*models.Record, error) {
	record, err := s.db.GetData(ctx, userName, meta)
	if err != nil {
		if errors.Is(err, errs.ErrRecordNotFound) {
			return nil, err
		}
		s.log.Errorw("failed to get credential", "user", userName, "meta", meta, "error", err)
		return nil, errs.ErrInternal
	}

	data, err := s.cipher.Decrypt(record.Data)
	if err != nil {
		s.log.Errorw("failed to decrypt data", "user", userName, "meta", meta, "error", err)
		return nil, errs.ErrInternal
	}
	record.Data = data

	return record, nil
}

// GetListMeta returns a list of user metadata records.
func (s *service) GetListMeta(ctx context.Context, userName string) ([]string, error) {
	records, err := s.db.GetListMeta(ctx, userName)
	if err != nil {
		s.log.Errorw("failed to get credentials list", "user", userName, "error", err)
		return nil, errs.ErrInternal
	}

	return records, nil
}

// UpdateData encrypts and updates existing user data.
func (s *service) UpdateData(ctx context.Context, record *models.Record) error {
	encrypted, err := s.cipher.Encrypt(record.Data)
	if err != nil {
		s.log.Errorw("failed to encrypt data", "user", record.UserName, "error", err)
		return errs.ErrInternal
	}

	if err := s.db.UpdateData(ctx, record.Type, record.Meta, record.UserName, encrypted); err != nil {
		if errors.Is(err, errs.ErrRecordNotFound) {
			return err
		}
		s.log.Errorw("failed to update credential", "error", err)
		return errs.ErrInternal
	}
	return nil
}

// DeleteData removes user data from storage.
func (s *service) DeleteData(ctx context.Context, userName, meta string) error {
	if err := s.db.DeleteData(ctx, userName, meta); err != nil {
		if errors.Is(err, errs.ErrRecordNotFound) {
			return err
		}
		s.log.Errorw("failed to delete credential", "user", userName, "meta", meta, "error", err)
		return errs.ErrInternal
	}
	return nil
}
