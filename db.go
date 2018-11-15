package main

import (
	"github.com/dgraph-io/badger"
	"github.com/pkg/errors"
)

var db *badger.DB

func InitDB(path string) error {
	opts := badger.DefaultOptions
	opts.Dir = *dbPath
	opts.ValueDir = *dbPath
	var err error
	if db, err = badger.Open(opts); err != nil {
		return errors.Wrap(err, "init db error")
	}
	return nil
}

func Get(key string) ([]byte, error) {
	var val []byte
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		val, err = item.ValueCopy(nil)
		return err
	})
	if err != nil {
		return nil, errors.Wrap(err, "get error")
	}

	return val, nil
}

// Set wrapper badger Set function.
func Set(key string, val []byte) error {
	if err := db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), val)
	}); err != nil {
		return errors.Wrap(err, "set error")
	}
	return nil
}

// Delete wrapper badger Delete function.
func Delete(key string) error {
	if err := db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	}); err != nil {
		return errors.Wrap(err, "delete error")
	}
	return nil
}

func ListKeys(prefix string) ([]string, error) {
	var keys []string
	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		p := []byte(prefix)
		for it.Seek(p); it.ValidForPrefix(p); it.Next() {
			k := it.Item().KeyCopy(nil)
			keys = append(keys, string(k))
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "list keys error")
	}
	return keys, nil
}
