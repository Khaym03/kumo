package storage

import (
	badger "github.com/dgraph-io/badger/v4"
)

var (
	htmlPrefix = []byte("html:")
	pdfPrefix  = []byte("pdf:")
	jsonPrefix = []byte("json:")
)

func (b *BadgerDBStore) SaveHTML(url string, htmlContent []byte) error {
	key := append(htmlPrefix, []byte(url)...)
	return b.save(key, htmlContent)
}

func (b *BadgerDBStore) SavePDF(id string, data []byte) error {
	key := append(pdfPrefix, []byte(id)...)
	return b.save(key, data)
}

func (b *BadgerDBStore) SaveJSON(id string, data []byte) error {
	key := append(jsonPrefix, []byte(id)...)
	return b.save(key, data)
}

func (b *BadgerDBStore) GetHTML(k string) ([]byte, error) {
	key := append(htmlPrefix, []byte(k)...)
	return b.get(key)
}

func (b *BadgerDBStore) GetPDF(k string) ([]byte, error) {
	key := append(pdfPrefix, []byte(k)...)
	return b.get(key)
}

func (b *BadgerDBStore) GetJSON(k string) ([]byte, error) {
	key := append(jsonPrefix, []byte(k)...)
	return b.get(key)
}

func (b *BadgerDBStore) get(key []byte) ([]byte, error) {
	var value []byte
	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		value, err = item.ValueCopy(nil)
		return err
	})
	return value, err
}

func (b *BadgerDBStore) save(key, data []byte) error {
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, data)
	})
}
