package flags

import (
	"fmt"
	"path/filepath"

	"github.com/lucasepe/locker/internal/kv"
	"github.com/lucasepe/locker/internal/kv/bbolt"
	"github.com/lucasepe/strcase"
)

const (
	defaultStoreName = "locker"
)

type Store struct {
	BaseDir      string
	MasterSecret string

	path string
	ref  kv.Store
}

func (f *Store) String() string {
	return f.path
}

func (f *Store) Set(v string) (err error) {
	name := v[:len(v)-len(filepath.Ext(v))]
	name = strcase.Kebab(name)
	f.path = filepath.Join(f.BaseDir, fmt.Sprintf("%s.db", name))

	return nil
}

func (f *Store) Connect() (kv.Store, error) {
	if f.ref != nil {
		return f.ref, nil
	}

	if len(f.path) == 0 {
		err := f.Set(defaultStoreName)
		if err != nil {
			return nil, err
		}
	}

	opts := bbolt.Options{Path: f.path}
	if len(f.MasterSecret) > 0 {
		opts.Codec = kv.NewCryptoCodec(f.MasterSecret)
	}

	return bbolt.NewStore(opts)
}
