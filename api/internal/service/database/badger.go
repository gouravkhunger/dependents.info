package database

import (
	"log"
	"strings"
	"time"

	"github.com/dgraph-io/badger/v4"
)

var (
	openAttempts = 120
	openWait     = time.Second
)

type Txn = *badger.Txn
type IteratorOptions = *badger.IteratorOptions

type BadgerService struct {
	db     *badger.DB
	gcStop chan struct{}
}

func NewBadgerService(path string) *BadgerService {
	opts := badger.DefaultOptions(path)
	opts = opts.WithLoggingLevel(badger.WARNING).
		WithNumLevelZeroTablesStall(10).
		WithValueLogFileSize(256 << 20).
		WithIndexCacheSize(64 << 20).
		WithBaseTableSize(16 << 20).
		WithValueThreshold(1 << 20).
		WithMemTableSize(32 << 20).
		WithNumLevelZeroTables(3)

	db, err := openDB(opts)
	if err != nil {
		panic("Failed to open Badger database: " + err.Error())
	}

	service := &BadgerService{
		db:     db,
		gcStop: make(chan struct{}),
	}

	go func() {
		ticker := time.NewTicker(6 * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				err := service.db.RunValueLogGC(0.7)
				if err != nil && err != badger.ErrNoRewrite {
					log.Printf("GC error: %v", err)
				}
			case <-service.gcStop:
				return
			}
		}
	}()

	return service
}

func openDB(opts badger.Options) (*badger.DB, error) {
	var db *badger.DB
	var err error
	for i := 0; i < openAttempts; i++ {
		db, err = badger.Open(opts)
		if err == nil {
			return db, nil
		}
		if !isLockErr(err) || i == openAttempts-1 {
			return nil, err
		}
		if i == 0 {
			log.Printf("waiting for badger lock")
		}
		time.Sleep(openWait)
	}
	return nil, err
}

func isLockErr(err error) bool {
	return err != nil && strings.Contains(err.Error(), "directory lock")
}

func (b *BadgerService) Sync() error {
	return b.db.Sync()
}

func (b *BadgerService) Close() error {
	close(b.gcStop)
	return b.db.Close()
}
