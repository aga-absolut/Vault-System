package models

type Record struct {
	ID       int64
	Type     string
	Meta     string
	Data     []byte
	UserName string
}

type User struct {
	Name     string
	Password string
}

type SecretConfig struct {
	TokenTTL      string `json:"token_ttl"`
	JWTSecret     string `json:"jwt_secret"`
	EncryptionKey string `json:"encryption_key"`
}
