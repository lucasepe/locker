package store_test

import (
	"bytes"
	"log"
	"os"
	"testing"
	"time"

	"github.com/lucasepe/locker/internal/store"
	"go.etcd.io/bbolt"
)

// Ensure that we can create and delete a bucket.
func TestBucket(t *testing.T) {
	db := newTestDB()
	defer db.Close()

	_, err := store.NewBucket(db, []byte("things"))
	if err != nil {
		t.Fatal(err)
	}

	err = store.DeleteBucket(db, []byte("things"))
	if err != nil {
		t.Fatal(err)
	}
}

// Ensure we can put an item in a bucket.
func TestPut(t *testing.T) {
	db := newTestDB()
	defer db.Close()

	things, err := store.NewBucket(db, []byte("things"))
	if err != nil {
		t.Fatal(err)
	}

	key, value := []byte("A"), []byte("alpha")

	// Put key/value into the `things` bucket.
	if err := things.Put(key, value); err != nil {
		t.Fatal(err)
	}

	// Read value back in a different read-only transaction.
	got, err := things.Get(key)
	if err != nil && !bytes.Equal(got, value) {
		t.Fatal(err)
	}
}

// Ensure we don't overwrite existing items when using PutNX.
func TestPutNX(t *testing.T) {
	db := newTestDB()
	defer db.Close()

	things, err := store.NewBucket(db, []byte("things"))
	if err != nil {
		t.Fatal(err)
	}

	key := []byte("A")
	a, b := []byte("alpha"), []byte("beta")

	// Put key/a-value into the `things` bucket.
	if err := things.PutNX(key, a); err != nil {
		t.Fatal(err)
	}

	// Read value back in a different read-only transaction.
	got, err := things.Get(key)
	if err != nil && !bytes.Equal(got, a) {
		t.Fatal(err)
	}

	// Try putting key/b-value into the `things` bucket.
	if err := things.PutNX(key, b); err != nil {
		t.Fatal(err)
	}

	// Value for key should still be a, not b.
	got, err = things.Get(key)
	if err != nil && !bytes.Equal(got, a) {
		t.Fatal(err)
	}
}

// Ensure that a bucket that gets a non-existent key returns nil.
func TestGetMissing(t *testing.T) {
	db := newTestDB()
	defer db.Close()

	things, err := store.NewBucket(db, []byte("things"))
	if err != nil {
		t.Fatal(err)
	}

	key := []byte("missing")
	if got, _ := things.Get(key); got != nil {
		t.Errorf("not expecting value for key %q: got %q", key, got)
	}
}

// Ensure that we can delete stuff in a bucket.
func TestDelete(t *testing.T) {
	db := newTestDB()
	defer db.Close()

	things, err := store.NewBucket(db, []byte("things"))
	if err != nil {
		t.Fatal(err)
	}

	k, v := []byte("foo"), []byte("bar")
	if err = things.Put(k, v); err != nil {
		t.Fatal(err)
	}

	if err = things.Delete(k); err != nil {
		t.Fatal(err)
	}
}

// Ensure we can insert items into a bucket and get them back out.
func TestInsert(t *testing.T) {
	db := newTestDB()
	defer db.Close()

	paths, err := store.NewBucket(db, []byte("paths"))
	if err != nil {
		t.Fatal(err)
	}

	// k, v pairs to put in `paths` bucket
	items := []struct {
		Key, Value []byte
	}{
		{[]byte("foo/"), []byte("foo")},
		{[]byte("foo/bar/"), []byte("bar")},
		{[]byte("foo/bar/baz/"), []byte("baz")},
		{[]byte("food/"), []byte("")},
		{[]byte("good/"), []byte("")},
		{[]byte("goo/"), []byte("")},
	}

	err = paths.Insert(items)
	if err != nil {
		t.Fatal(err)
	}

	gotItems, err := paths.Items()
	if err != nil {
		t.Fatal(err)
	}

	// expected k/v mapping
	expected := map[string][]byte{
		"foo/":         []byte("foo"),
		"foo/bar/":     []byte("bar"),
		"foo/bar/baz/": []byte("baz"),
		"food/":        []byte(""),
		"good/":        []byte(""),
		"goo/":         []byte(""),
	}

	for _, item := range gotItems {
		want := expected[string(item.Key)]
		if !bytes.Equal(item.Value, want) {
			t.Fatalf("got %v, want %v", item.Value, want)
		}
	}
}

// Ensure we can safely insert items into a bucket without overwriting
// existing items.
func TestInsertNX(t *testing.T) {
	db := newTestDB()
	defer db.Close()

	bk, err := store.NewBucket(db, []byte("test"))
	if err != nil {
		t.Fatal(err)
	}

	// Put k/v into the `bk` bucket.
	k, v := []byte("A"), []byte("alpha")
	if err := bk.Put(k, v); err != nil {
		t.Fatal(err)
	}

	// k/v pairs to put-if-not-exists
	items := []struct {
		Key, Value []byte
	}{
		{[]byte("A"), []byte("ALPHA")}, // key exists, so don't update
		{[]byte("B"), []byte("beta")},
		{[]byte("C"), []byte("gamma")},
	}

	err = bk.InsertNX(items)
	if err != nil {
		t.Fatal(err)
	}

	gotItems, err := bk.Items()
	if err != nil {
		t.Fatal(err)
	}

	// expected items
	expected := []struct {
		Key, Value []byte
	}{
		{[]byte("A"), []byte("alpha")}, // existing value not updated
		{[]byte("B"), []byte("beta")},
		{[]byte("C"), []byte("gamma")},
	}

	for i, got := range gotItems {
		want := expected[i]
		if !bytes.Equal(got.Value, want.Value) {
			t.Fatalf("key %q: got %v, want %v", got.Key, got.Value, want.Value)
		}
	}
}

// newTestDB returns a TestDB using a temporary path.
func newTestDB() *bbolt.DB {
	db, err := bbolt.Open(tempfile(), 0600, &bbolt.Options{
		Timeout: 30 * time.Second,
	})
	if err != nil {
		log.Fatalf("cannot open buckets database: %s", err)
	}
	return db
}

// tempfile returns a temporary file path.
func tempfile() string {
	f, err := os.CreateTemp("", "bolt-")
	if err != nil {
		log.Fatalf("Could not create temp file: %s", err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
	if err := os.Remove(f.Name()); err != nil {
		log.Fatal(err)
	}
	return f.Name()
}
