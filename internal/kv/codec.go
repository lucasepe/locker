package kv

import (
	"encoding/base64"
	"os"

	"github.com/lucasepe/locker/internal/secrets"
)

// Codec encodes/decodes secrets to/from slices of bytes.
type Codec interface {
	// Marshal encodes a secret to a slice of bytes.
	Marshal(sec []byte) ([]byte, error)
	// Unmarshal decodes a slice of bytes into a secret.
	Unmarshal(data []byte) ([]byte, error)
}

func NewCryptoCodec() Codec {
	return &cryptoCodec{}
}

var _ Codec = (*cryptoCodec)(nil)

type cryptoCodec struct{}

func (cc *cryptoCodec) Marshal(src []byte) ([]byte, error) {
	token := os.Getenv(EnvSecret)
	if len(token) == 0 {
		return nil, ErrUnsetMasterPassword
	}

	dat, err := secrets.Encrypt([]byte(token), src)
	if err != nil {
		return nil, err
	}

	enc := base64.StdEncoding
	buf := make([]byte, enc.EncodedLen(len(dat)))
	enc.Encode(buf, dat)

	return buf, nil
}

func (cc *cryptoCodec) Unmarshal(data []byte) ([]byte, error) {
	token := os.Getenv(EnvSecret)
	if len(token) == 0 {
		return nil, ErrUnsetMasterPassword
	}

	enc := base64.StdEncoding
	dbuf := make([]byte, enc.DecodedLen(len(data)))
	n, err := enc.Decode(dbuf, data)
	if err != nil {
		return nil, err
	}

	return secrets.Decrypt([]byte(token), dbuf[:n])
}
