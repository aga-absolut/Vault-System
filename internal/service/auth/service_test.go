package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/aga-absolut/Vault-System/internal/errs"
	"github.com/aga-absolut/Vault-System/internal/interceptors/logger"
	mockDB "github.com/aga-absolut/Vault-System/internal/storage/postgres/mocks"
	mockToken "github.com/aga-absolut/Vault-System/internal/token/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_RegisterUser(t *testing.T) {
	tests := []struct {
		name          string
		username      string
		password      string
		setupmockDB   func(*mockDB.MockStorage, *mockToken.MockProvider)
		expectedToken string
		expectedErr   error
	}{
		{
			name:     "succesful registration",
			username: "newuser",
			password: "strongpass123",
			setupmockDB: func(db *mockDB.MockStorage, token *mockToken.MockProvider) {
				db.EXPECT().AddUser(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
				token.EXPECT().BuildJWTString(gomock.Any()).Return("jwt-token-for-newuser-42", nil)
			},
			expectedToken: "jwt-token-for-newuser-42",
			expectedErr:   nil,
		},
		{
			name:     "login already used",
			username: "existing",
			password: "password123",
			setupmockDB: func(db *mockDB.MockStorage, _ *mockToken.MockProvider) {
				db.EXPECT().AddUser(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errs.ErrLoginAlreadyUsed)
			},
			expectedToken: "",
			expectedErr:   errs.ErrLoginAlreadyUsed,
		},
		{
			name:     "empty login",
			username: "",
			password: "pass123",
			setupmockDB: func(db *mockDB.MockStorage, _ *mockToken.MockProvider) {
				// ничего не вызываем
			},
			expectedToken: "",
			expectedErr:   errs.ErrIncorrectLoginOrPassword,
		},
		{
			name:     "too short password",
			username: "validuser",
			password: "123",
			setupmockDB: func(db *mockDB.MockStorage, _ *mockToken.MockProvider) {
				// ничего не вызываем
			},
			expectedToken: "",
			expectedErr:   errs.ErrTooShortPassword,
		},
		{
			name:     "failed create user",
			username: "newuser2",
			password: "strongpass123",
			setupmockDB: func(db *mockDB.MockStorage, _ *mockToken.MockProvider) {
				db.EXPECT().AddUser(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errs.ErrInternal)
			},
			expectedToken: "",
			expectedErr:   errs.ErrInternal,
		},
		{
			name:     "failed build token",
			username: "newuser3",
			password: "strongpass123",
			setupmockDB: func(db *mockDB.MockStorage, token *mockToken.MockProvider) {
				db.EXPECT().AddUser(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
				token.EXPECT().BuildJWTString(gomock.Any()).Return("", errs.ErrInternal)
			},
			expectedToken: "",
			expectedErr:   errs.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockdb := mockDB.NewMockStorage(ctrl)
			mockToken := mockToken.NewMockProvider(ctrl)
			tt.setupmockDB(mockdb, mockToken)

			log := logger.NewLogger()
			svc := NewService(log, mockdb, mockToken)

			tokenStr, err := svc.RegisterUser(context.Background(), tt.username, tt.password)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.True(t, errors.Is(err, tt.expectedErr), "ожидали %v, получили %v", tt.expectedErr, err)
			} else {
				require.NoError(t, err)
			}

			if tt.expectedToken != "" {
				assert.NotEmpty(t, tokenStr)
			} else {
				assert.Empty(t, tokenStr)
			}
		})
	}
}

func TestAuthService_LoginUser(t *testing.T) {
	var (
		validUser = "validuser"
		validPass = "correctpassword123"
	)

	validHashBytes, _ := bcrypt.GenerateFromPassword([]byte(validPass), bcrypt.DefaultCost)
	validHash := string(validHashBytes)

	tests := []struct {
		name          string
		username      string
		password      string
		setupmockDB   func(*mockDB.MockStorage, *mockToken.MockProvider)
		expectedToken string
		expectedErr   error
	}{
		{
			name:     "succesful authorization",
			username: validUser,
			password: validPass,
			setupmockDB: func(db *mockDB.MockStorage, token *mockToken.MockProvider) {
				db.EXPECT().CheckUser(gomock.Any(), gomock.Any()).
					Return(validHash, nil)

				token.EXPECT().BuildJWTString(gomock.Any()).Return("jwt-token-success-12345", nil)
			},
			expectedToken: "jwt-token-success-12345",
			expectedErr:   nil,
		},
		{
			name:     "wrong password",
			username: "existing",
			password: "password123",
			setupmockDB: func(db *mockDB.MockStorage, token *mockToken.MockProvider) {
				db.EXPECT().CheckUser(gomock.Any(), gomock.Any()).
					Return("", errs.ErrIncorrectLoginOrPassword)
			},
			expectedToken: "",
			expectedErr:   errs.ErrIncorrectLoginOrPassword,
		},
		{
			name:     "failed authorization from db",
			username: "newuser2",
			password: "strongpass123",
			setupmockDB: func(db *mockDB.MockStorage, _ *mockToken.MockProvider) {
				db.EXPECT().CheckUser(gomock.Any(), gomock.Any()).
					Return("", errs.ErrInternal)
			},
			expectedToken: "",
			expectedErr:   errs.ErrInternal,
		},
		{
			name:     "failed authorization from create token",
			username: validUser,
			password: validPass,
			setupmockDB: func(db *mockDB.MockStorage, token *mockToken.MockProvider) {
				db.EXPECT().CheckUser(gomock.Any(), gomock.Any()).
					Return(validHash, nil)

				token.EXPECT().BuildJWTString(gomock.Any()).Return("", errs.ErrInternal)
			},
			expectedToken: "",
			expectedErr:   errs.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDB := mockDB.NewMockStorage(ctrl)
			mockToken := mockToken.NewMockProvider(ctrl)
			tt.setupmockDB(mockDB, mockToken)

			log := logger.NewLogger()

			svc := NewService(log, mockDB, mockToken)

			tokenStr, err := svc.LoginUser(context.Background(), tt.username, tt.password)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.True(t, errors.Is(err, tt.expectedErr), "ожидали %v, получили %v", tt.expectedErr, err)
			} else {
				require.NoError(t, err)
			}

			if tt.expectedToken != "" {
				assert.NotEmpty(t, tokenStr)
			} else {
				assert.Empty(t, tokenStr)
			}
		})
	}
}
