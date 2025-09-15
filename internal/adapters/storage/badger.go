package storage

import (
	"encoding/json"

	"github.com/Khaym03/kumo/internal/pkg/types"
	badger "github.com/dgraph-io/badger/v4"
	log "github.com/sirupsen/logrus"
)

var (
	pendingKeyPrefix   = []byte("pending:")
	completedKeyPrefix = []byte("completed:")
)

type BadgerDBStore struct {
	db *badger.DB
}

func NewBadgerDB(dbPath string, allowLogger bool) (*badger.DB, error) {
	opts := badger.DefaultOptions(dbPath)
	opts.ValueLogFileSize = 32 << 20 // 32MB
	if !allowLogger {
		opts.Logger = nil

	}

	return badger.Open(opts)
}

func NewBadgerDBStore(db *badger.DB) *BadgerDBStore {
	return &BadgerDBStore{db: db}
}

func (b *BadgerDBStore) SavePending(requests ...*types.Request) error {
	txn := b.db.NewTransaction(true)
	defer txn.Discard()

	for _, req := range requests {
		key := append(pendingKeyPrefix, []byte(req.URL)...)
		val, err := json.Marshal(req)
		if err != nil {
			return err
		}
		if err := txn.Set(key, val); err != nil {
			return err
		}
	}
	return txn.Commit()
}

func (b *BadgerDBStore) LoadPending() ([]*types.Request, error) {
	var requests []*types.Request
	err := b.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := pendingKeyPrefix
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			val, err := item.ValueCopy(nil)
			if err != nil {
				log.Warnf("Error copying value: %v", err)
				continue
			}
			var req types.Request
			if err := json.Unmarshal(val, &req); err != nil {
				log.Warnf("Error unmarshalling request: %v", err)
				continue
			}
			requests = append(requests, &req)
		}
		return nil
	})
	return requests, err
}

func (b *BadgerDBStore) SaveCompleted(req *types.Request) error {
	return b.db.Update(func(txn *badger.Txn) error {
		completedKey := append(completedKeyPrefix, []byte(req.URL)...)
		val, err := json.Marshal(req)
		if err != nil {
			return err
		}
		return txn.Set(completedKey, val)
	})
}

func (b *BadgerDBStore) RemoveFromPending(req *types.Request) error {
	return b.db.Update(func(txn *badger.Txn) error {
		pendingKey := append(pendingKeyPrefix, []byte(req.URL)...)
		return txn.Delete(pendingKey)
	})
}

func (b *BadgerDBStore) Close() error {
	return b.db.Close()
}

func (b *BadgerDBStore) IsCompleted(url string) (bool, error) {
	key := append(completedKeyPrefix, []byte(url)...)
	var found bool
	err := b.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get(key)
		if err == nil {
			found = true
			return nil
		}
		if err == badger.ErrKeyNotFound {
			found = false
			return nil
		}
		return err
	})
	return found, err
}
