package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
)

type Cipher struct {
	key []byte
}

func NewCipher(keyStr string) *Cipher {
	keyStr = strings.TrimSpace(keyStr)
	decoded, err := base64.StdEncoding.DecodeString(keyStr)
	if err != nil {
		decoded = []byte(keyStr)
	}

	if len(decoded) != 32 {
		panic(fmt.Sprintf("AES key must be 32 bytes after decoding, got %d", len(decoded)))
	}

	return &Cipher{key: decoded}
}

func (c *Cipher) Encrypt(data []byte) ([]byte, error) {
	aesblock, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, err
	}

	nonce, err := generateNonce(aesgcm.NonceSize())
	if err != nil {
		return nil, err
	}

	ciphertext := aesgcm.Seal(nil, nonce, data, nil)
	combined := append(nonce, ciphertext...)

	return combined, nil
}

func (c *Cipher) Decrypt(combined []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := combined[:aesgcm.NonceSize()]
	ciphertext := combined[aesgcm.NonceSize():]

	data, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func generateNonce(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
