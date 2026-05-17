package credentials

import (
	"context"
	"testing"

	"github.com/aga-absolut/Vault-System/internal/config"
	"github.com/aga-absolut/Vault-System/internal/crypto"
	"github.com/aga-absolut/Vault-System/internal/errs"
	"github.com/aga-absolut/Vault-System/internal/interceptors/logger"
	"github.com/aga-absolut/Vault-System/internal/models"
	"github.com/aga-absolut/Vault-System/internal/storage/postgres/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCredentialsService_SetData(t *testing.T) {
	tests := []struct {
		name        string
		record      *models.Record
		setupMockDB func(*mocks.MockStorage)
		expectedErr error
	}{
		{
			name: "succesful set data",
			record: &models.Record{
				Type:     "email",
				Meta:     "absolute@email.ru",
				Data:     []byte("password123"),
				UserName: "absolute_1",
			},
			setupMockDB: func(db *mocks.MockStorage) {
				db.EXPECT().SetData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "meta already used",
			record: &models.Record{
				ID:       67,
				Type:     "email",
				Meta:     "absolute@email.ru",
				Data:     []byte("password123"),
				UserName: "absolute_1",
			},
			setupMockDB: func(db *mocks.MockStorage) {
				db.EXPECT().SetData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errs.ErrMetaAlreadyUsed)
			},
			expectedErr: errs.ErrMetaAlreadyUsed,
		},
		{
			name: "failed set data",
			record: &models.Record{
				ID:       67,
				Type:     "email",
				Meta:     "absolute@email.ru",
				Data:     []byte("password123"),
				UserName: "absolute_1",
			},
			setupMockDB: func(db *mocks.MockStorage) {
				db.EXPECT().SetData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errs.ErrInternal)
			},
			expectedErr: errs.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks.NewMockStorage(ctrl)
			tt.setupMockDB(m)

			cfg := &config.Config{
				EncryptionKey: "32-bytes-long-secret-key-exactly",
			}
			log := logger.NewLogger()
			cipher := crypto.NewCipher(cfg.EncryptionKey)
			srv := NewService(log, m, cipher)

			err := srv.SetData(context.Background(), tt.record)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCredentialsService_GetData(t *testing.T) {
	cipher := crypto.NewCipher("32-bytes-long-secret-key-exactly")
	encryptedData, _ := cipher.Encrypt([]byte("password123"))
	fakeEncryptedData := []byte("wrongPassword")

	tests := []struct {
		name        string
		userName    string
		meta        string
		setupMockDB func(*mocks.MockStorage)
		wantData    []byte
		expectedErr error
	}{
		{
			name:     "succesful get data",
			userName: "absolute_1",
			meta:     "abs",
			setupMockDB: func(db *mocks.MockStorage) {
				db.EXPECT().GetData(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&models.Record{Data: encryptedData}, nil)
			},
			wantData:    []byte("password123"),
			expectedErr: nil,
		},
		{
			name:     "error record not found",
			userName: "absolute_1",
			meta:     "abs",
			setupMockDB: func(db *mocks.MockStorage) {
				db.EXPECT().GetData(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&models.Record{}, errs.ErrRecordNotFound)
			},
			wantData:    nil,
			expectedErr: errs.ErrRecordNotFound,
		},
		{
			name:     "error server internal",
			userName: "absolute_1",
			meta:     "abs",
			setupMockDB: func(db *mocks.MockStorage) {
				db.EXPECT().GetData(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&models.Record{}, errs.ErrInternal)
			},
			wantData:    nil,
			expectedErr: errs.ErrInternal,
		},
		{
			name:     "failed to decrypt record",
			userName: "absolute_1",
			meta:     "abs",
			setupMockDB: func(db *mocks.MockStorage) {
				db.EXPECT().GetData(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&models.Record{Data: fakeEncryptedData}, errs.ErrInternal)
			},
			wantData:    nil,
			expectedErr: errs.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks.NewMockStorage(ctrl)
			tt.setupMockDB(m)

			cfg := &config.Config{
				EncryptionKey: "32-bytes-long-secret-key-exactly",
			}
			log := logger.NewLogger()
			cipher := crypto.NewCipher(cfg.EncryptionKey)
			srv := NewService(log, m, cipher)

			record, err := srv.GetData(context.Background(), tt.userName, tt.meta)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}

			if tt.wantData != nil {
				assert.Equal(t, record.Data, tt.wantData)
			}
		})
	}
}

func TestCredentialsService_GetListData(t *testing.T) {
	listMeta := []string{"aaa", "bbb", "ccc"}

	tests := []struct {
		name         string
		userName     string
		setupMockDB  func(*mocks.MockStorage)
		wantListData []string
		expectedErr  error
	}{
		{
			name:     "succesful get data",
			userName: "absolute_1",
			setupMockDB: func(db *mocks.MockStorage) {
				db.EXPECT().GetListMeta(gomock.Any(), gomock.Any()).
					Return(listMeta, nil)
			},
			wantListData: listMeta,
			expectedErr:  nil,
		},
		{
			name:     "error server internal",
			userName: "absolute_1",
			setupMockDB: func(db *mocks.MockStorage) {
				db.EXPECT().GetListMeta(gomock.Any(), gomock.Any()).
					Return(nil, errs.ErrInternal)
			},
			wantListData: nil,
			expectedErr:  errs.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks.NewMockStorage(ctrl)
			tt.setupMockDB(m)

			cfg := &config.Config{
				EncryptionKey: "32-bytes-long-secret-key-exactly",
			}
			log := logger.NewLogger()
			cipher := crypto.NewCipher(cfg.EncryptionKey)
			srv := NewService(log, m, cipher)

			list, err := srv.GetListMeta(context.Background(), tt.userName)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}

			if tt.wantListData != nil {
				for i, meta := range list {
					assert.Equal(t, tt.wantListData[i], meta)
				}
			}
		})
	}
}

func TestCredentialsService_UpdateData(t *testing.T) {
	tests := []struct {
		name        string
		record      *models.Record
		setupMockDB func(*mocks.MockStorage)
		expectedErr error
	}{
		{
			name: "succesful set data",
			record: &models.Record{
				Type:     "email",
				Meta:     "absolute@email.ru",
				Data:     []byte("password123"),
				UserName: "absolute_1",
			},
			setupMockDB: func(db *mocks.MockStorage) {
				db.EXPECT().UpdateData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "meta already used",
			record: &models.Record{
				Type:     "email",
				Meta:     "absolute@email.ru",
				Data:     []byte("password123"),
				UserName: "absolute_1",
			},
			setupMockDB: func(db *mocks.MockStorage) {
				db.EXPECT().UpdateData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errs.ErrRecordNotFound)
			},
			expectedErr: errs.ErrRecordNotFound,
		},
		{
			name: "failed set data",
			record: &models.Record{
				Type:     "email",
				Meta:     "absolute@email.ru",
				Data:     []byte("password123"),
				UserName: "absolute_1",
			},
			setupMockDB: func(db *mocks.MockStorage) {
				db.EXPECT().UpdateData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errs.ErrInternal)
			},
			expectedErr: errs.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks.NewMockStorage(ctrl)
			tt.setupMockDB(m)

			cfg := &config.Config{
				EncryptionKey: "32-bytes-long-secret-key-exactly",
			}
			log := logger.NewLogger()
			cipher := crypto.NewCipher(cfg.EncryptionKey)
			srv := NewService(log, m, cipher)

			err := srv.UpdateData(context.Background(), tt.record)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCredentialsService_DeleteData(t *testing.T) {
	tests := []struct {
		name        string
		userName    string
		meta        string
		setupMockDB func(*mocks.MockStorage)
		expectedErr error
	}{
		{
			name:     "succesful get data",
			userName: "absolute_1",
			meta:     "aaa",
			setupMockDB: func(db *mocks.MockStorage) {
				db.EXPECT().DeleteData(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:     "error server internal",
			userName: "absolute_2",
			meta:     "aaa",
			setupMockDB: func(db *mocks.MockStorage) {
				db.EXPECT().DeleteData(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errs.ErrInternal)
			},
			expectedErr: errs.ErrInternal,
		},
		{
			name:     "error id not found",
			userName: "absolute_3",
			meta:     "aaa",
			setupMockDB: func(db *mocks.MockStorage) {
				db.EXPECT().DeleteData(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errs.ErrRecordNotFound)
			},
			expectedErr: errs.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks.NewMockStorage(ctrl)
			tt.setupMockDB(m)

			cfg := &config.Config{
				EncryptionKey: "32-bytes-long-secret-key-exactly",
			}
			log := logger.NewLogger()
			cipher := crypto.NewCipher(cfg.EncryptionKey)
			srv := NewService(log, m, cipher)

			err := srv.DeleteData(context.Background(), tt.userName, tt.meta)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
