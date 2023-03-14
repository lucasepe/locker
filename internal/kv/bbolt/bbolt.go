package bbolt

import (
	"github.com/lucasepe/locker/internal/kv"
	"go.etcd.io/bbolt"
)

// Options are the options for the bbolt store.
type Options struct {
	// Path of the DB file.
	// Optional ("bbolt.db" by default).
	Path string
	// Encoding format.
	// Optional (encoding.JSON by default).
	Codec kv.Codec
}

// DefaultOptions is an Options object with default values.
// Path: "bbolt.db", Codec: encoding.JSON
var DefaultOptions = Options{
	Path:  "secrets.db",
	Codec: kv.NewCryptoCodec(),
}

// NewStore creates a new bbolt store.
// Note: bbolt uses an exclusive write lock on the database file so it cannot be shared by multiple processes.
// So when creating multiple clients you should always use a new database file (by setting a different Path in the options).
//
// You must call the Close() method on the store when you're done working with it.
func NewStore(options Options) (kv.Store, error) {
	if options.Path == "" {
		options.Path = DefaultOptions.Path
	}
	if options.Codec == nil {
		options.Codec = kv.NewCryptoCodec()
	}

	// Open DB
	db, err := bbolt.Open(options.Path, 0600, nil)
	if err != nil {
		return nil, err
	}

	return &boltStore{
		db:    db,
		codec: options.Codec,
	}, nil
}

var _ kv.Store = (*boltStore)(nil)

// boltStore is a kv.Store implementation for bbolt (formerly known as Bolt / Bolt DB).
type boltStore struct {
	db    *bbolt.DB
	codec kv.Codec
}

func (s *boltStore) PutOne(namespace string, key, value string) error {
	if len(namespace) == 0 {
		return kv.ErrEmptyNamespace
	}

	if len(key) == 0 {
		return kv.ErrEmptyKey
	}

	data, err := s.codec.Marshal([]byte(value))
	if err != nil {
		return err
	}

	return s.db.Update(func(tx *bbolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists([]byte(namespace))
		if err != nil {
			return err
		}
		return bkt.Put([]byte(key), data)
	})
}

func (s *boltStore) GetOne(namespace, key string) (value string, err error) {
	if len(namespace) == 0 {
		return "", kv.ErrEmptyNamespace
	}

	if len(key) == 0 {
		return "", kv.ErrEmptyKey
	}

	err = s.db.View(func(tx *bbolt.Tx) error {
		bkt := tx.Bucket([]byte(namespace))
		data := bkt.Get([]byte(key))
		if data != nil {
			dst, err := s.codec.Unmarshal(data)
			if err != nil {
				return err
			}
			value = string(dst)
		}
		return nil
	})

	return value, err
}

func (s *boltStore) DeleteOne(namespace, key string) error {
	if len(namespace) == 0 {
		return kv.ErrEmptyNamespace
	}

	if len(key) == 0 {
		return kv.ErrEmptyKey
	}

	return s.db.Update(func(tx *bbolt.Tx) error {
		bkt := tx.Bucket([]byte(namespace))
		return bkt.Delete([]byte(key))
	})
}

func (s *boltStore) DeleteAll(namespace string) error {
	if len(namespace) == 0 {
		return kv.ErrEmptyNamespace
	}

	return s.db.Update(func(tx *bbolt.Tx) error {
		return tx.DeleteBucket([]byte(namespace))
	})
}

func (s *boltStore) GetAll(namespace string) (map[string]string, error) {
	res := make(map[string]string)

	err := s.db.View(func(tx *bbolt.Tx) error {
		bkt := tx.Bucket([]byte(namespace))
		if bkt == nil {
			return kv.ErrNamespaceNotFound
		}

		c := bkt.Cursor()
		for key, val := c.First(); key != nil; key, val = c.Next() {
			k := make([]byte, len(key))
			copy(k, key)

			v, err := s.codec.Unmarshal(val)
			if err != nil {
				return err
			}

			res[string(k)] = string(v)
		}

		return nil
	})

	return res, err
}

func (s *boltStore) Namespaces() (names []string, err error) {
	err = s.db.View(func(tx *bbolt.Tx) error {
		return tx.ForEach(func(bn []byte, _ *bbolt.Bucket) error {
			el := make([]byte, len(bn))
			copy(el, bn)
			names = append(names, string(el))
			return nil
		})
	})

	return names, err
}

// Keys returns all keys in a namespace.
func (s *boltStore) Keys(namespace string) (items []string, err error) {
	return items, s.db.View(func(tx *bbolt.Tx) error {
		bkt := tx.Bucket([]byte(namespace))
		if bkt == nil {
			return bbolt.ErrBucketNotFound
		}

		c := bkt.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			key := make([]byte, len(k))
			copy(key, k)

			items = append(items, string(key))
		}
		return nil
	})
}

// Close closes the store.
func (s *boltStore) Close() error {
	if s.db == nil {
		return nil
	}
	return s.db.Close()
}
