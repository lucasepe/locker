package app

import (
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

	return secrets.Encrypt(key, data)
}

func Decrypt(data []byte) ([]byte, error) {
	key := []byte(os.Getenv(EnvSecret))
	if len(key) == 0 {
		return nil, ErrUnsetMasterPassword
	}

	return secrets.Decrypt(key, data)
	//return nil, fmt.Errorf("%w: double check your master password", err)
}
