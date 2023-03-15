package flags

import "github.com/lucasepe/strcase"

type Namespace struct {
	name []byte
}

func (f *Namespace) String() string {
	return string(f.name)
}

func (f *Namespace) Set(v string) (err error) {
	f.name = []byte(strcase.Kebab(v))
	return nil
}

func (f *Namespace) Bytes() []byte {
	return f.name
}
