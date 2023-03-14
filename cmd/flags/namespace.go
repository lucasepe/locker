package flags

import "github.com/lucasepe/strcase"

type NamespaceFlag struct {
	name []byte
}

func (f *NamespaceFlag) String() string {
	return string(f.name)
}

func (f *NamespaceFlag) Set(v string) (err error) {
	f.name = []byte(strcase.Kebab(v))
	return nil
}

func (f *NamespaceFlag) Bytes() []byte {
	return f.name
}
