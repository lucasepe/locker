package app

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"github.com/lucasepe/locker/internal/secrets"
	"github.com/lucasepe/xdg"
)

const (
	Name      = "locker"
	EnvSecret = "LOCKER_SECRET"
)

var (
	ErrUnsetMasterPassword = fmt.Errorf(
		"specify a master password setting the env var: %s", EnvSecret)
)

func Dir() string {
	return filepath.Join(xdg.ConfigDir(), Name)
}

func MasterPassword() string {
	return os.Getenv(EnvSecret)
}

func Encrypt(data []byte) ([]byte, error) {
	key := []byte(os.Getenv(EnvSecret))
	if len(key) == 0 {
		return nil, ErrUnsetMasterPassword
	}

	src, err := secrets.Encrypt(key, data)
	if err != nil {
		return nil, err
	}

	encoder := base64.StdEncoding
	buf := make([]byte, encoder.EncodedLen(len(src)))
	encoder.Encode(buf, src)
	return buf, nil
}

func Decrypt(data []byte) ([]byte, error) {
	key := []byte(os.Getenv(EnvSecret))
	if len(key) == 0 {
		return nil, ErrUnsetMasterPassword
	}

	enc := base64.StdEncoding
	dbuf := make([]byte, enc.DecodedLen(len(data)))
	n, err := enc.Decode(dbuf, data)
	if err != nil {
		return nil, err
	}

	res, err := secrets.Decrypt(key, dbuf[:n])
	if err != nil {
		return nil, fmt.Errorf("%w: double check your master password", err)
	}

	return res, nil
}
