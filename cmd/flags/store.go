package flags

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lucasepe/locker/cmd/app"
	"github.com/lucasepe/locker/internal/kv"
	"github.com/lucasepe/locker/internal/kv/bbolt"
	"github.com/lucasepe/strcase"
)

type StoreFlag struct {
	path string
	name string
	ref  kv.Store
}

func (f *StoreFlag) String() string {
	return f.path
}

func (f *StoreFlag) Name() string {
	return f.name
}

func (f *StoreFlag) Set(name string) (err error) {
	f.name = strcase.Kebab(name)
	f.path = filepath.Join(app.Dir(), fmt.Sprintf("%s.db", f.name))

	return os.MkdirAll(app.Dir(), os.ModePerm)
}

func (f *StoreFlag) Connect() (kv.Store, error) {
	if f.ref != nil {
		return f.ref, nil
	}

	if len(f.path) == 0 {
		err := f.Set(app.Name)
		if err != nil {
			return nil, err
		}
	}

	return bbolt.NewStore(bbolt.Options{
		Path: f.path,
	})
}
