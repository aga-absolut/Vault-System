package postgres

import (
	"context"
	"errors"

	"github.com/aga-absolut/Vault-System/internal/config"
	"github.com/aga-absolut/Vault-System/internal/errs"
	"github.com/aga-absolut/Vault-System/internal/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

//go:generate mockgen -source=postgres.go -destination=mocks/postgres_mock.go -package=mocks
type Storage interface {
	AddUser(ctx context.Context, name, password string) error
	CheckUser(ctx context.Context, name string) (string, error)
	SetData(ctx context.Context, recType, recName, userName string, data []byte) error
	GetData(ctx context.Context, userName, meta string) (*models.Record, error)
	UpdateData(ctx context.Context, recType, recName, userName string, data []byte) error
	GetListMeta(ctx context.Context, userName string) ([]string, error)
	DeleteData(ctx context.Context, userName, meta string) error
	Close()
}

type storage struct {
	DB     *pgxpool.Pool
	config *config.Config
}

func NewStorage(config *config.Config) Storage {
	pgx, err := pgxpool.New(context.Background(), config.DatabaseDSN)
	if err != nil {
		panic(err)
	}

	return &storage{
		DB:     pgx,
		config: config,
	}
}

func (s *storage) AddUser(ctx context.Context, name, password string) error {
	query := `INSERT INTO users (username, password_hash) VALUES($1, $2)`

	if _, err := s.DB.Exec(ctx, query, name, password); err != nil {
		var PgErr *pgconn.PgError
		if errors.As(err, &PgErr) && PgErr.Code == pgerrcode.UniqueViolation {
			return errs.ErrLoginAlreadyUsed
		}
		return err
	}
	return nil
}

func (s *storage) CheckUser(ctx context.Context, name string) (string, error) {
	var passwordHash string
	query := `SELECT password_hash FROM users WHERE username = $1`

	if err := s.DB.QueryRow(ctx, query, name).Scan(&passwordHash); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errs.ErrIncorrectLoginOrPassword
		}
		return "", err
	}

	return passwordHash, nil
}

func (s *storage) SetData(ctx context.Context, recType, recMeta, userName string, data []byte) error {
	query := `INSERT INTO records (record_type, record_meta, record_data, username) 
	          VALUES($1, $2, $3, $4)`

	if _, err := s.DB.Exec(ctx, query, recType, recMeta, data, userName); err != nil {
		var PgErr *pgconn.PgError
		if errors.As(err, &PgErr) && PgErr.Code == pgerrcode.UniqueViolation {
			return errs.ErrMetaAlreadyUsed
		}
		return err
	}
	return nil
}

func (s *storage) GetData(ctx context.Context, userName, meta string) (*models.Record, error) {
	var r models.Record
	query := `SELECT record_id, record_type, record_meta, record_data, username FROM records 
	WHERE record_meta = $1 AND username = $2`

	err := s.DB.QueryRow(ctx, query, meta, userName).Scan(&r.ID, &r.Type, &r.Meta, &r.Data, &r.UserName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrRecordNotFound
		}
		return nil, err
	}
	return &r, nil
}

func (s *storage) UpdateData(ctx context.Context, recType, recMeta, userName string, data []byte) error {
	query := `UPDATE records SET record_type = $1, record_data = $2 WHERE record_meta = $3 AND username = $4`
	result, err := s.DB.Exec(ctx, query, recType, data, recMeta, userName)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errs.ErrRecordNotFound
	}

	return nil
}

func (s *storage) GetListMeta(ctx context.Context, userName string) ([]string, error) {
	query := `SELECT record_meta FROM records WHERE username = $1`
	rows, err := s.DB.Query(ctx, query, userName)
	if err != nil {
		return nil, err
	}

	var records []string
	for rows.Next() {
		var record string
		if err := rows.Scan(&record); err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

func (s *storage) DeleteData(ctx context.Context, userName, meta string) error {
	query := `DELETE FROM records WHERE record_meta = $1 AND username = $2`
	result, err := s.DB.Exec(ctx, query, meta, userName)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errs.ErrRecordNotFound
	}
	return nil
}

func (s *storage) Close() {
	s.DB.Close()
}
