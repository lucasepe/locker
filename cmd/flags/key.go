package flags

import "github.com/lucasepe/strcase"

type Key struct {
	content []byte
}

func (f *Key) String() string {
	return string(f.content)
}

func (f *Key) Set(v string) (err error) {
	f.content = []byte(strcase.Snake(v))
	return nil
}

func (f *Key) Bytes() []byte {
	return f.content
}
