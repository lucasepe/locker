package flags

import (
	"fmt"
)

type StringList struct {
	array []string
}

func (f *StringList) String() string {
	return fmt.Sprint(f.array)
}

func (f *StringList) Set(v string) error {
	if f.array == nil {
		f.array = make([]string, 1)
	} else {
		nv := make([]string, len(f.array)+1)
		copy(nv, f.array)
		f.array = nv
	}

	f.array[len(f.array)-1] = v

	return nil
}

func (f *StringList) Values() []string {
	return f.array
}
