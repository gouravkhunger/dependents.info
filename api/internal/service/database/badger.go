package database

import (
	"github.com/dgraph-io/badger/v4"
)

type Txn = *badger.Txn

type BadgerService struct {
	db *badger.DB
}

func NewBadgerService(path string) *BadgerService {
	opts := badger.DefaultOptions(path)
	opts = opts.WithLoggingLevel(badger.WARNING)

	db, err := badger.Open(opts)
	if err != nil {
		panic("Failed to open Badger database: " + err.Error())
	}

	return &BadgerService{db: db}
}

func (b *BadgerService) Sync() error {
	return b.db.Sync()
}

func (b *BadgerService) Close() error {
	return b.db.Close()
}
