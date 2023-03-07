package flags

import "path/filepath"

type FileFlag struct {
	path string
}

func (f *FileFlag) String() string {
	return f.path
}

func (f *FileFlag) Set(v string) (err error) {
	f.path, err = filepath.Abs(v)
	if err != nil {
		return err
	}
	return nil
}
