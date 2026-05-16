package models

type (
	// Record represents stored user credential data.
	Record struct {
		ID       int64
		Type     string
		Meta     string
		Data     []byte
		UserName string
	}

	// User represents application user credentials.
	User struct {
		Name     string
		Password string
	}

	// SecretConfig stores sensitive application configuration values.
	SecretConfig struct {
		TokenTTL      string `json:"token_ttl"`
		JWTSecret     string `json:"jwt_secret"`
		EncryptionKey string `json:"encryption_key"`
	}
)
