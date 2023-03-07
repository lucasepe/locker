package secrets

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestSimple(t *testing.T) {
	key := []byte("hello world")
	data := []byte("hello jello")

	encdata, err := Encrypt(key, data)
	if err != nil {
		t.Fatal(err)
	}

	decdata, err := Decrypt(key, encdata)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(decdata, data) {
		t.Fatalf("expected: %s, got: %s",
			hex.EncodeToString(data), hex.EncodeToString(decdata))
	}
}
