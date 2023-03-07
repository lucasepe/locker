package flags

type BoxFlag struct {
	name []byte
}

func (f *BoxFlag) String() string {
	return string(f.name)
}

func (f *BoxFlag) Set(v string) (err error) {
	f.name = []byte(v)
	return nil
}

func (f *BoxFlag) Bytes() []byte {
	return f.name
}
