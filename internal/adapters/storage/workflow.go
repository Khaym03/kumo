package storage

import (
	"encoding/json"

	"github.com/Khaym03/kumo/internal/ports"
	badger "github.com/dgraph-io/badger/v4"
	log "github.com/sirupsen/logrus"
)

var (
	workFlowPrefix  = []byte("workflow:")
)

type BadgerWorkFlowStore struct {
	db *badger.DB
}

func NewBadgerWorkFlowStore(db *badger.DB) ports.WorkFlow {
	return  &BadgerWorkFlowStore{db: db}
}


func (b *BadgerWorkFlowStore) Save(k string, v map[string]any) error {
	txn := b.db.NewTransaction(true)
	defer txn.Discard()
	defer txn.Commit()

	value, err:=json.Marshal(v)
	if err != nil {
		return  err
	}
	key := append(workFlowPrefix, []byte(k)...)

	log.Infof("Original:%v Copy:%v", string(workFlowPrefix), string(key))

	return txn.Set(key, value)
}

func (b *BadgerWorkFlowStore) Load() ([]map[string]any, error) {
	var requests []map[string]any
	err := b.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := workFlowPrefix
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			log.Println("item: ", item)
			val, err := item.ValueCopy(nil)
			if err != nil {
				log.Warnf("Error copying value: %v", err)
				continue
			}
			var req map[string]any
			if err := json.Unmarshal(val, &req); err != nil {
				log.Warnf("Error unmarshalling request: %v", err)
				continue
			}
			requests = append(requests, req)
		}
		return nil
	})
	return requests, err
}