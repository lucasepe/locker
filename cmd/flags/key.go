package flags

import "github.com/lucasepe/strcase"

type KeyFlag struct {
	content []byte
}

func (f *KeyFlag) String() string {
	return string(f.content)
}

func (f *KeyFlag) Set(v string) (err error) {
	f.content = []byte(strcase.Snake(v))
	return nil
}

func (f *KeyFlag) Bytes() []byte {
	return f.content
}
