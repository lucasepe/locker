package kv

import (
	"fmt"
)

func ExampleCodec_Marshal() {
	codec := NewCryptoCodec("MAGIK")
	enc, err := codec.Marshal([]byte("Hello World!"))
	if err != nil {
		panic(err)
	}

	dec, err := codec.Unmarshal(enc)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(dec))

	// Output:
	// Hello World!
}
