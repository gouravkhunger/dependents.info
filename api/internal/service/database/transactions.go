package database

import (
	"github.com/dgraph-io/badger/v4"
)

var opts IteratorOptions

func init() {
	opts = &badger.DefaultIteratorOptions
	opts.PrefetchValues = false
}

func (b *BadgerService) Get(key string, out *string) error {
	var data []byte
	err := b.db.View(func(txn Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			data = make([]byte, len(val))
			copy(data, val)
			return nil
		})
	})

	if err != nil {
		return err
	}

	*out = string(data)
	return nil
}

func (b *BadgerService) Save(key string, data []byte) error {
	return b.db.Update(func(txn Txn) error {
		return txn.Set([]byte(key), data)
	})
}

func (b *BadgerService) IterateKeys(callback func(key string)) {
	b.db.View(func(txn Txn) error {
		it := txn.NewIterator(*opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			key := it.Item().KeyCopy(nil)
			callback(string(key))
		}
		return nil
	})
}
