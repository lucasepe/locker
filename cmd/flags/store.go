package flags

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/lucasepe/locker/cmd/app"
	"github.com/lucasepe/strcase"
	"go.etcd.io/bbolt"
)

type StoreFlag struct {
	path string
	name string
	ref  *bbolt.DB
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

func (f *StoreFlag) Connect() (*bbolt.DB, error) {
	if f.ref != nil {
		return f.ref, nil
	}

	if len(f.path) == 0 {
		err := f.Set(app.Name)
		if err != nil {
			return nil, err
		}
	}

	return bbolt.Open(f.path, 0600, &bbolt.Options{
		Timeout: 30 * time.Second,
	})
}
