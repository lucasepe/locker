package flags

import "github.com/lucasepe/strcase"

type LabelFlag struct {
	content []byte
}

func (f *LabelFlag) String() string {
	return string(f.content)
}

func (f *LabelFlag) Set(v string) (err error) {
	f.content = []byte(strcase.Camel(v))
	return nil
}

func (f *LabelFlag) Bytes() []byte {
	return f.content
}
